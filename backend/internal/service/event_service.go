package service

import (
	"errors"
	"ticketeer/backend/internal/domain"
	"ticketeer/backend/internal/repository"

	"gorm.io/gorm"
)

type EventService struct {
	eventRepository *repository.EventRepository
}

func NewEventService(eventRepository *repository.EventRepository) *EventService {
	return &EventService{
		eventRepository: eventRepository,
	}
}

func (s *EventService) GetEvents() ([]domain.Event, error) {
	return s.eventRepository.FindAll()
}

func (s *EventService) GetEvent(id uint) (*domain.Event, error) {
	event, err := s.eventRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return event, nil
}
