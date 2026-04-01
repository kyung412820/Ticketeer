package tests

import (
	"net/http"
	"sync"
	"sync/atomic"
	"testing"

	"ticketeer/backend/internal/domain"
	"ticketeer/backend/internal/testutil"
)

func TestConcurrentHoldOnlyOneSucceeds(t *testing.T) {
	db, rdb, app, err := testutil.SetupTestApp()
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := testutil.ResetState(db, rdb); err != nil {
		t.Fatalf("reset: %v", err)
	}

	event, _ := testutil.SeedOpenEvent(db)
	seat, _ := testutil.SeedSeat(db, event.ID, "A-20", domain.SeatStatusAvailable)

	const workers = 10

	type actor struct {
		clientID   string
		queueToken string
	}

	actors := make([]actor, 0, workers)
	for i := 0; i < workers; i++ {
		clientID := "client-hold-" + string(rune('a'+i))
		enter := performJSON(t, app, http.MethodPost, "/api/queue/enter", map[string]any{
			"event_id":  event.ID,
			"client_id": clientID,
		})
		if enter.Code != http.StatusOK {
			t.Fatalf("queue enter failed for %s: %s", clientID, enter.Body.String())
		}
		queueToken := parseJSON(t, enter)["data"].(map[string]any)["queue_token"].(string)
		actors = append(actors, actor{
			clientID:   clientID,
			queueToken: queueToken,
		})
	}

	var successCount int32
	var failureCount int32

	var wg sync.WaitGroup
	wg.Add(workers)

	for _, a := range actors {
		a := a
		go func() {
			defer wg.Done()

			hold := performJSON(t, app, http.MethodPost, "/api/seats/1/hold", map[string]any{
				"event_id":    event.ID,
				"client_id":   a.clientID,
				"queue_token": a.queueToken,
			})

			if hold.Code == http.StatusOK {
				atomic.AddInt32(&successCount, 1)
				return
			}

			if hold.Code == http.StatusBadRequest || hold.Code == http.StatusNotFound {
				atomic.AddInt32(&failureCount, 1)
				return
			}

			t.Errorf("unexpected status code: %d body=%s", hold.Code, hold.Body.String())
		}()
	}

	wg.Wait()

	if successCount != 1 {
		t.Fatalf("expected exactly 1 hold success, got %d", successCount)
	}

	if failureCount != workers-1 {
		t.Fatalf("expected %d hold failures, got %d", workers-1, failureCount)
	}

	seats := performJSON(t, app, http.MethodGet, "/api/events/1/seats", nil)
	if seats.Code != http.StatusOK {
		t.Fatalf("seat fetch failed: %s", seats.Body.String())
	}

	body := parseJSON(t, seats)
	data := body["data"].(map[string]any)
	seatList := data["seats"].([]any)
	if len(seatList) != 1 {
		t.Fatalf("expected 1 seat, got %d", len(seatList))
	}

	seatMap := seatList[0].(map[string]any)
	if seatMap["status"] != "HELD" {
		t.Fatalf("expected seat status HELD after concurrent holds, got %v", seatMap["status"])
	}

	_ = seat
}
