package tests

import (
	"context"
	"fmt"
)

func ctx() context.Context {
	return context.Background()
}

func holdKey(eventID, seatID uint) string {
	return fmt.Sprintf("seat_hold:%d:%d", eventID, seatID)
}
