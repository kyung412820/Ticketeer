package service

import "errors"

var (
	ErrEventNotFound      = errors.New("event not found")
	ErrBookingNotOpen     = errors.New("booking not open")
	ErrEventClosed        = errors.New("event closed")
	ErrAlreadyInQueue     = errors.New("already in queue")
	ErrQueueTokenNotFound = errors.New("queue token not found")
	ErrInvalidQueueToken  = errors.New("invalid queue token")
	ErrQueueExpired       = errors.New("queue expired")
	ErrSeatNotFound       = errors.New("seat not found")
	ErrSeatEventMismatch  = errors.New("seat event mismatch")
	ErrSeatAlreadyBooked  = errors.New("seat already booked")
	ErrSeatAlreadyHeld    = errors.New("seat already held")
	ErrHoldNotFound       = errors.New("hold not found")
	ErrHoldExpired        = errors.New("hold expired")
)
