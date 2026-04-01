package tests

import (
	"net/http"
	"sync"
	"sync/atomic"
	"testing"

	"ticketeer/backend/internal/domain"
	"ticketeer/backend/internal/testutil"
)

func TestConcurrentBookingOnlyOneSucceeds(t *testing.T) {
	db, rdb, app, err := testutil.SetupTestApp()
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := testutil.ResetState(db, rdb); err != nil {
		t.Fatalf("reset: %v", err)
	}

	event, _ := testutil.SeedOpenEvent(db)
	seat, _ := testutil.SeedSeat(db, event.ID, "A-10", domain.SeatStatusAvailable)

	enter := performJSON(t, app, http.MethodPost, "/api/queue/enter", map[string]any{
		"event_id":  event.ID,
		"client_id": "client-concurrent",
	})
	if enter.Code != http.StatusOK {
		t.Fatalf("queue enter failed: %s", enter.Body.String())
	}

	queueToken := parseJSON(t, enter)["data"].(map[string]any)["queue_token"].(string)

	hold := performJSON(t, app, http.MethodPost, "/api/seats/1/hold", map[string]any{
		"event_id":    event.ID,
		"client_id":   "client-concurrent",
		"queue_token": queueToken,
	})
	if hold.Code != http.StatusOK {
		t.Fatalf("hold failed: %s", hold.Body.String())
	}

	const workers = 10

	var successCount int32
	var failureCount int32

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			book := performJSON(t, app, http.MethodPost, "/api/bookings", map[string]any{
				"event_id":    event.ID,
				"seat_id":     seat.ID,
				"client_id":   "client-concurrent",
				"queue_token": queueToken,
			})

			if book.Code == http.StatusOK {
				atomic.AddInt32(&successCount, 1)
				return
			}

			if book.Code == http.StatusBadRequest || book.Code == http.StatusNotFound {
				atomic.AddInt32(&failureCount, 1)
				return
			}

			t.Errorf("unexpected status code: %d body=%s", book.Code, book.Body.String())
		}()
	}

	wg.Wait()

	if successCount != 1 {
		t.Fatalf("expected exactly 1 success, got %d", successCount)
	}

	if failureCount != workers-1 {
		t.Fatalf("expected %d failures, got %d", workers-1, failureCount)
	}

	var bookingCount int64
	if err := db.Model(&domain.Booking{}).Where("seat_id = ?", seat.ID).Count(&bookingCount).Error; err != nil {
		t.Fatalf("count bookings: %v", err)
	}

	if bookingCount != 1 {
		t.Fatalf("expected exactly 1 booking row, got %d", bookingCount)
	}

	var updatedSeat domain.Seat
	if err := db.First(&updatedSeat, seat.ID).Error; err != nil {
		t.Fatalf("reload seat: %v", err)
	}

	if updatedSeat.Status != domain.SeatStatusBooked {
		t.Fatalf("expected seat status BOOKED, got %s", updatedSeat.Status)
	}
}
