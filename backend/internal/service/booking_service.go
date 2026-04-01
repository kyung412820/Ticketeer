package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ticketeer/backend/internal/domain"
	"ticketeer/backend/internal/repository"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BookingResponse struct {
	BookingID   uint      `json:"booking_id"`
	BookingCode string    `json:"booking_code"`
	EventID     uint      `json:"event_id"`
	SeatID      uint      `json:"seat_id"`
	Status      string    `json:"status"`
	BookedAt    time.Time `json:"booked_at"`
}

type BookingService struct {
	db                *gorm.DB
	seatRepository    *repository.SeatRepository
	queueRepository   *repository.QueueRepository
	bookingRepository *repository.BookingRepository
	rdb               *redis.Client
}

func NewBookingService(
	db *gorm.DB,
	seatRepository *repository.SeatRepository,
	queueRepository *repository.QueueRepository,
	bookingRepository *repository.BookingRepository,
	rdb *redis.Client,
) *BookingService {
	return &BookingService{
		db:                db,
		seatRepository:    seatRepository,
		queueRepository:   queueRepository,
		bookingRepository: bookingRepository,
		rdb:               rdb,
	}
}

func (s *BookingService) CreateBooking(eventID, seatID uint, clientID, queueToken string) (*BookingResponse, error) {
	entry, err := s.queueRepository.FindByToken(queueToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQueueTokenNotFound
		}
		return nil, err
	}

	if entry.EventID != eventID || entry.ClientID != clientID {
		return nil, ErrInvalidQueueToken
	}

	now := time.Now()
	if entry.Status == domain.QueueStatusReady && entry.ExpiredAt != nil && now.After(*entry.ExpiredAt) {
		entry.Status = domain.QueueStatusExpired
		if err := s.queueRepository.Save(entry); err != nil {
			return nil, err
		}
		return nil, ErrQueueExpired
	}

	if entry.Status != domain.QueueStatusReady {
		return nil, ErrInvalidQueueToken
	}

	ctx := context.Background()
	holdKey := fmt.Sprintf("seat_hold:%d:%d", eventID, seatID)

	holdClientID, err := s.rdb.Get(ctx, holdKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrHoldNotFound
		}
		return nil, err
	}
	if holdClientID != clientID {
		return nil, ErrHoldExpired
	}

	bookingCode := "BK-" + uuid.NewString()
	bookedAt := time.Now()
	var createdBooking domain.Booking

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		lockedSeat, err := s.seatRepository.FindByIDForUpdate(tx, seatID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrSeatNotFound
			}
			return err
		}

		if lockedSeat.EventID != eventID {
			return ErrSeatEventMismatch
		}

		if lockedSeat.Status == domain.SeatStatusBooked {
			return ErrSeatAlreadyBooked
		}

		lockedSeat.Status = domain.SeatStatusBooked
		if err := s.seatRepository.SaveTx(tx, lockedSeat); err != nil {
			return err
		}

		createdBooking = domain.Booking{
			EventID:     eventID,
			SeatID:      seatID,
			ClientID:    clientID,
			BookingCode: bookingCode,
			Status:      domain.BookingStatusConfirmed,
			BookedAt:    bookedAt,
		}
		if err := s.bookingRepository.CreateTx(tx, &createdBooking); err != nil {
			return err
		}

		entry.Status = domain.QueueStatusExpired
		if err := s.queueRepository.SaveTx(tx, entry); err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		switch {
		case errors.Is(txErr, ErrSeatNotFound):
			return nil, ErrSeatNotFound
		case errors.Is(txErr, ErrSeatEventMismatch):
			return nil, ErrSeatEventMismatch
		case errors.Is(txErr, ErrSeatAlreadyBooked):
			return nil, ErrSeatAlreadyBooked
		default:
			return nil, txErr
		}
	}

	if err := s.rdb.Del(ctx, holdKey).Err(); err != nil {
		return nil, err
	}

	return &BookingResponse{
		BookingID:   createdBooking.ID,
		BookingCode: createdBooking.BookingCode,
		EventID:     createdBooking.EventID,
		SeatID:      createdBooking.SeatID,
		Status:      string(createdBooking.Status),
		BookedAt:    createdBooking.BookedAt,
	}, nil
}
