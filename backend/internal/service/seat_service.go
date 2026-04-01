package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ticketeer/backend/internal/domain"
	"ticketeer/backend/internal/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const seatHoldTTL = 5 * time.Minute

type SeatResponse struct {
	ID      uint   `json:"id"`
	SeatNo  string `json:"seat_no"`
	Section string `json:"section"`
	Price   int    `json:"price"`
	Status  string `json:"status"`
}

type HoldSeatResponse struct {
	SeatID        uint      `json:"seat_id"`
	Status        string    `json:"status"`
	HoldExpiresAt time.Time `json:"hold_expires_at"`
}

type SeatService struct {
	seatRepository  *repository.SeatRepository
	queueRepository *repository.QueueRepository
	rdb             *redis.Client
}

func NewSeatService(seatRepository *repository.SeatRepository, queueRepository *repository.QueueRepository, rdb *redis.Client) *SeatService {
	return &SeatService{
		seatRepository:  seatRepository,
		queueRepository: queueRepository,
		rdb:             rdb,
	}
}

func (s *SeatService) GetSeatsByEventID(eventID uint) ([]SeatResponse, error) {
	seats, err := s.seatRepository.FindByEventID(eventID)
	if err != nil {
		return nil, err
	}

	result := make([]SeatResponse, 0, len(seats))
	ctx := context.Background()

	for _, seat := range seats {
		status := string(seat.Status)

		if seat.Status == domain.SeatStatusAvailable {
			holdKey := fmt.Sprintf("seat_hold:%d:%d", eventID, seat.ID)
			exists, err := s.rdb.Exists(ctx, holdKey).Result()
			if err == nil && exists > 0 {
				status = string(domain.SeatStatusHeld)
			}
		}

		result = append(result, SeatResponse{
			ID:      seat.ID,
			SeatNo:  seat.SeatNo,
			Section: seat.Section,
			Price:   seat.Price,
			Status:  status,
		})
	}

	return result, nil
}

func (s *SeatService) HoldSeat(seatID uint, eventID uint, clientID, queueToken string) (*HoldSeatResponse, error) {
	seat, err := s.seatRepository.FindByID(seatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSeatNotFound
		}
		return nil, err
	}

	if seat.EventID != eventID {
		return nil, ErrSeatEventMismatch
	}

	if seat.Status == domain.SeatStatusBooked {
		return nil, ErrSeatAlreadyBooked
	}

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
	expiresAt := now.Add(seatHoldTTL)

	ok, err := s.rdb.SetNX(ctx, holdKey, clientID, seatHoldTTL).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrSeatAlreadyHeld
	}

	return &HoldSeatResponse{
		SeatID:        seatID,
		Status:        string(domain.SeatStatusHeld),
		HoldExpiresAt: expiresAt,
	}, nil
}
