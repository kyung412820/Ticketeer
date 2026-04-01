package repository

import (
	"ticketeer/backend/internal/domain"

	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(booking *domain.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepository) CreateTx(tx *gorm.DB, booking *domain.Booking) error {
	return tx.Create(booking).Error
}
