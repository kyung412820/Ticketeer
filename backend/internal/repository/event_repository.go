package repository

import (
	"ticketeer/backend/internal/domain"

	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) FindAll() ([]domain.Event, error) {
	var events []domain.Event
	if err := r.db.Order("event_at asc").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *EventRepository) FindByID(id uint) (*domain.Event, error) {
	var event domain.Event
	if err := r.db.First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}
