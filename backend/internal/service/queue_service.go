package service

import (
	"errors"
	"time"

	"ticketeer/backend/internal/domain"
	"ticketeer/backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	queueReadyWindow = 10 * time.Minute
	queueReadyLimit  = 5
)

type QueueService struct {
	eventRepository *repository.EventRepository
	queueRepository *repository.QueueRepository
}

type QueueEnterResponse struct {
	QueueToken string `json:"queue_token"`
	Status     string `json:"status"`
	Position   int    `json:"position"`
}

type QueueStatusResponse struct {
	QueueToken string     `json:"queue_token"`
	Status     string     `json:"status"`
	Position   int        `json:"position,omitempty"`
	ReadyAt    *time.Time `json:"ready_at,omitempty"`
	ExpiredAt  *time.Time `json:"expired_at,omitempty"`
}

func NewQueueService(eventRepository *repository.EventRepository, queueRepository *repository.QueueRepository) *QueueService {
	return &QueueService{
		eventRepository: eventRepository,
		queueRepository: queueRepository,
	}
}

func (s *QueueService) Enter(eventID uint, clientID string) (*QueueEnterResponse, error) {
	event, err := s.eventRepository.FindByID(eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	now := time.Now()
	if now.Before(event.BookingOpenAt) || event.Status == domain.EventStatusOpenPending {
		return nil, ErrBookingNotOpen
	}
	if now.After(event.BookingCloseAt) || event.Status == domain.EventStatusClosed {
		return nil, ErrEventClosed
	}

	active, err := s.queueRepository.FindActiveByEventAndClient(eventID, clientID)
	if err == nil && active != nil {
		return nil, ErrAlreadyInQueue
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	waitingCount, err := s.queueRepository.CountWaiting(eventID)
	if err != nil {
		return nil, err
	}

	entry := &domain.QueueEntry{
		EventID:    eventID,
		ClientID:   clientID,
		QueueToken: uuid.NewString(),
		Status:     domain.QueueStatusWaiting,
		Position:   int(waitingCount) + 1,
	}

	if entry.Position <= queueReadyLimit {
		readyAt := now
		expiredAt := now.Add(queueReadyWindow)
		entry.Status = domain.QueueStatusReady
		entry.Position = 0
		entry.ReadyAt = &readyAt
		entry.ExpiredAt = &expiredAt
	}

	if err := s.queueRepository.Create(entry); err != nil {
		return nil, err
	}

	return &QueueEnterResponse{
		QueueToken: entry.QueueToken,
		Status:     string(entry.Status),
		Position:   entry.Position,
	}, nil
}

func (s *QueueService) GetStatus(token string) (*QueueStatusResponse, error) {
	entry, err := s.queueRepository.FindByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQueueTokenNotFound
		}
		return nil, err
	}

	now := time.Now()

	if entry.Status == domain.QueueStatusReady && entry.ExpiredAt != nil && now.After(*entry.ExpiredAt) {
		entry.Status = domain.QueueStatusExpired
		if err := s.queueRepository.Save(entry); err != nil {
			return nil, err
		}
	}

	if entry.Status == domain.QueueStatusWaiting {
		count, err := s.queueRepository.CountWaitingBefore(entry.EventID, entry.CreatedAt)
		if err == nil {
			entry.Position = int(count)
			if entry.Position <= queueReadyLimit {
				readyAt := now
				expiredAt := now.Add(queueReadyWindow)
				entry.Status = domain.QueueStatusReady
				entry.Position = 0
				entry.ReadyAt = &readyAt
				entry.ExpiredAt = &expiredAt
				if err := s.queueRepository.Save(entry); err != nil {
					return nil, err
				}
			}
		}
	}

	resp := &QueueStatusResponse{
		QueueToken: entry.QueueToken,
		Status:     string(entry.Status),
		Position:   entry.Position,
		ReadyAt:    entry.ReadyAt,
		ExpiredAt:  entry.ExpiredAt,
	}
	return resp, nil
}
