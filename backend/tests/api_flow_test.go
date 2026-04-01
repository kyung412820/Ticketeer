package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ticketeer/backend/internal/domain"
	"ticketeer/backend/internal/testutil"
)

func performJSON(t *testing.T, handler http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func parseJSON(t *testing.T, rr *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("parse json: %v, body=%s", err, rr.Body.String())
	}
	return out
}

func TestHappyPath_BookingFlow(t *testing.T) {
	db, rdb, app, err := testutil.SetupTestApp()
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := testutil.ResetState(db, rdb); err != nil {
		t.Fatalf("reset: %v", err)
	}

	event, _ := testutil.SeedOpenEvent(db)
	seat, _ := testutil.SeedSeat(db, event.ID, "A-1", domain.SeatStatusAvailable)

	enter := performJSON(t, app, http.MethodPost, "/api/queue/enter", map[string]any{
		"event_id":  event.ID,
		"client_id": "client-a",
	})
	if enter.Code != http.StatusOK {
		t.Fatalf("queue enter failed: %s", enter.Body.String())
	}
	queueToken := parseJSON(t, enter)["data"].(map[string]any)["queue_token"].(string)

	hold := performJSON(t, app, http.MethodPost, "/api/seats/1/hold", map[string]any{
		"event_id":    event.ID,
		"client_id":   "client-a",
		"queue_token": queueToken,
	})
	if hold.Code != http.StatusOK {
		t.Fatalf("hold failed: %s", hold.Body.String())
	}

	book := performJSON(t, app, http.MethodPost, "/api/bookings", map[string]any{
		"event_id":    event.ID,
		"seat_id":     seat.ID,
		"client_id":   "client-a",
		"queue_token": queueToken,
	})
	if book.Code != http.StatusOK {
		t.Fatalf("booking failed: %s", book.Body.String())
	}

	bookData := parseJSON(t, book)["data"].(map[string]any)
	if bookData["status"] != "CONFIRMED" {
		t.Fatalf("expected CONFIRMED, got %v", bookData["status"])
	}
}

func TestBookingNotOpen(t *testing.T) {
	db, rdb, app, err := testutil.SetupTestApp()
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := testutil.ResetState(db, rdb); err != nil {
		t.Fatalf("reset: %v", err)
	}

	event, _ := testutil.SeedPendingEvent(db)

	enter := performJSON(t, app, http.MethodPost, "/api/queue/enter", map[string]any{
		"event_id":  event.ID,
		"client_id": "client-a",
	})
	if enter.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", enter.Code, enter.Body.String())
	}

	errMap := parseJSON(t, enter)["error"].(map[string]any)
	if errMap["code"] != "BOOKING_NOT_OPEN" {
		t.Fatalf("expected BOOKING_NOT_OPEN, got %v", errMap["code"])
	}
}

func TestSeatAlreadyBooked(t *testing.T) {
	db, rdb, app, err := testutil.SetupTestApp()
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := testutil.ResetState(db, rdb); err != nil {
		t.Fatalf("reset: %v", err)
	}

	event, _ := testutil.SeedOpenEvent(db)
	_, _ = testutil.SeedSeat(db, event.ID, "A-1", domain.SeatStatusBooked)

	enter := performJSON(t, app, http.MethodPost, "/api/queue/enter", map[string]any{
		"event_id":  event.ID,
		"client_id": "client-a",
	})
	queueToken := parseJSON(t, enter)["data"].(map[string]any)["queue_token"].(string)

	hold := performJSON(t, app, http.MethodPost, "/api/seats/1/hold", map[string]any{
		"event_id":    event.ID,
		"client_id":   "client-a",
		"queue_token": queueToken,
	})
	if hold.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", hold.Code, hold.Body.String())
	}

	body := parseJSON(t, hold)
	if body["error"].(map[string]any)["code"] != "SEAT_ALREADY_BOOKED" {
		t.Fatalf("unexpected error: %v", body)
	}
}

func TestBookingWithoutHold(t *testing.T) {
	db, rdb, app, err := testutil.SetupTestApp()
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := testutil.ResetState(db, rdb); err != nil {
		t.Fatalf("reset: %v", err)
	}

	event, _ := testutil.SeedOpenEvent(db)
	seat, _ := testutil.SeedSeat(db, event.ID, "A-2", domain.SeatStatusAvailable)

	enter := performJSON(t, app, http.MethodPost, "/api/queue/enter", map[string]any{
		"event_id":  event.ID,
		"client_id": "client-b",
	})
	queueToken := parseJSON(t, enter)["data"].(map[string]any)["queue_token"].(string)

	book := performJSON(t, app, http.MethodPost, "/api/bookings", map[string]any{
		"event_id":    event.ID,
		"seat_id":     seat.ID,
		"client_id":   "client-b",
		"queue_token": queueToken,
	})
	if book.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", book.Code, book.Body.String())
	}

	body := parseJSON(t, book)
	if body["error"].(map[string]any)["code"] != "HOLD_NOT_FOUND" {
		t.Fatalf("unexpected error: %v", body)
	}
}

func TestInvalidQueueTokenOnHold(t *testing.T) {
	db, rdb, app, err := testutil.SetupTestApp()
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := testutil.ResetState(db, rdb); err != nil {
		t.Fatalf("reset: %v", err)
	}

	event, _ := testutil.SeedOpenEvent(db)
	_, _ = testutil.SeedSeat(db, event.ID, "A-3", domain.SeatStatusAvailable)

	hold := performJSON(t, app, http.MethodPost, "/api/seats/1/hold", map[string]any{
		"event_id":    event.ID,
		"client_id":   "client-c",
		"queue_token": "invalid-token",
	})
	if hold.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", hold.Code, hold.Body.String())
	}

	body := parseJSON(t, hold)
	if body["error"].(map[string]any)["code"] != "QUEUE_TOKEN_NOT_FOUND" {
		t.Fatalf("unexpected error: %v", body)
	}
}

func TestSeatAlreadyHeld(t *testing.T) {
	db, rdb, app, err := testutil.SetupTestApp()
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := testutil.ResetState(db, rdb); err != nil {
		t.Fatalf("reset: %v", err)
	}

	event, _ := testutil.SeedOpenEvent(db)
	_, _ = testutil.SeedSeat(db, event.ID, "A-4", domain.SeatStatusAvailable)

	enter1 := performJSON(t, app, http.MethodPost, "/api/queue/enter", map[string]any{
		"event_id":  event.ID,
		"client_id": "client-d",
	})
	token1 := parseJSON(t, enter1)["data"].(map[string]any)["queue_token"].(string)

	enter2 := performJSON(t, app, http.MethodPost, "/api/queue/enter", map[string]any{
		"event_id":  event.ID,
		"client_id": "client-e",
	})
	token2 := parseJSON(t, enter2)["data"].(map[string]any)["queue_token"].(string)

	hold1 := performJSON(t, app, http.MethodPost, "/api/seats/1/hold", map[string]any{
		"event_id":    event.ID,
		"client_id":   "client-d",
		"queue_token": token1,
	})
	if hold1.Code != http.StatusOK {
		t.Fatalf("first hold failed: %s", hold1.Body.String())
	}

	hold2 := performJSON(t, app, http.MethodPost, "/api/seats/1/hold", map[string]any{
		"event_id":    event.ID,
		"client_id":   "client-e",
		"queue_token": token2,
	})
	if hold2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", hold2.Code, hold2.Body.String())
	}

	body := parseJSON(t, hold2)
	if body["error"].(map[string]any)["code"] != "SEAT_ALREADY_HELD" {
		t.Fatalf("unexpected error: %v", body)
	}
}

func TestHoldMissingOnBooking(t *testing.T) {
	db, rdb, app, err := testutil.SetupTestApp()
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := testutil.ResetState(db, rdb); err != nil {
		t.Fatalf("reset: %v", err)
	}

	event, _ := testutil.SeedOpenEvent(db)
	seat, _ := testutil.SeedSeat(db, event.ID, "A-5", domain.SeatStatusAvailable)

	enter := performJSON(t, app, http.MethodPost, "/api/queue/enter", map[string]any{
		"event_id":  event.ID,
		"client_id": "client-f",
	})
	token := parseJSON(t, enter)["data"].(map[string]any)["queue_token"].(string)

	hold := performJSON(t, app, http.MethodPost, "/api/seats/1/hold", map[string]any{
		"event_id":    event.ID,
		"client_id":   "client-f",
		"queue_token": token,
	})
	if hold.Code != http.StatusOK {
		t.Fatalf("hold failed: %s", hold.Body.String())
	}

	if err := rdb.Del(ctx(), holdKey(event.ID, seat.ID)).Err(); err != nil {
		t.Fatalf("delete hold key: %v", err)
	}

	book := performJSON(t, app, http.MethodPost, "/api/bookings", map[string]any{
		"event_id":    event.ID,
		"seat_id":     seat.ID,
		"client_id":   "client-f",
		"queue_token": token,
	})
	if book.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", book.Code, book.Body.String())
	}

	body := parseJSON(t, book)
	if body["error"].(map[string]any)["code"] != "HOLD_NOT_FOUND" {
		t.Fatalf("unexpected error: %v", body)
	}
}
