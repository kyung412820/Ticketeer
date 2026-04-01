package repository

import (
	"time"

	"ticketeer/backend/internal/domain"

	"gorm.io/gorm"
)

type QueueRepository struct {
	db *gorm.DB
}

func NewQueueRepository(db *gorm.DB) *QueueRepository {
	return &QueueRepository{db: db}
}

func (r *QueueRepository) CountWaitingBefore(eventID uint, createdAt time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&domain.QueueEntry{}).
		Where("event_id = ? AND status = ? AND created_at <= ?", eventID, domain.QueueStatusWaiting, createdAt).
		Count(&count).Error
	return count, err
}

func (r *QueueRepository) CountWaiting(eventID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.QueueEntry{}).
		Where("event_id = ? AND status = ?", eventID, domain.QueueStatusWaiting).
		Count(&count).Error
	return count, err
}

func (r *QueueRepository) FindActiveByEventAndClient(eventID uint, clientID string) (*domain.QueueEntry, error) {
	var entry domain.QueueEntry
	err := r.db.
		Where("event_id = ? AND client_id = ? AND status IN ?", eventID, clientID, []domain.QueueStatus{domain.QueueStatusWaiting, domain.QueueStatusReady}).
		Order("created_at desc").
		First(&entry).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *QueueRepository) Create(entry *domain.QueueEntry) error {
	return r.db.Create(entry).Error
}

func (r *QueueRepository) FindByToken(token string) (*domain.QueueEntry, error) {
	var entry domain.QueueEntry
	if err := r.db.Where("queue_token = ?", token).First(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *QueueRepository) Save(entry *domain.QueueEntry) error {
	return r.db.Save(entry).Error
}

func (r *QueueRepository) SaveTx(tx *gorm.DB, entry *domain.QueueEntry) error {
	return tx.Save(entry).Error
}
