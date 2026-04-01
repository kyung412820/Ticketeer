package repository

import (
	"ticketeer/backend/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SeatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) *SeatRepository {
	return &SeatRepository{db: db}
}

func (r *SeatRepository) FindByEventID(eventID uint) ([]domain.Seat, error) {
	var seats []domain.Seat
	if err := r.db.
		Where("event_id = ?", eventID).
		Order("section asc, seat_no asc").
		Find(&seats).Error; err != nil {
		return nil, err
	}
	return seats, nil
}

func (r *SeatRepository) FindByID(id uint) (*domain.Seat, error) {
	var seat domain.Seat
	if err := r.db.First(&seat, id).Error; err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *SeatRepository) FindByIDForUpdate(tx *gorm.DB, id uint) (*domain.Seat, error) {
	var seat domain.Seat
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&seat, id).Error; err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *SeatRepository) Save(seat *domain.Seat) error {
	return r.db.Save(seat).Error
}

func (r *SeatRepository) SaveTx(tx *gorm.DB, seat *domain.Seat) error {
	return tx.Save(seat).Error
}
