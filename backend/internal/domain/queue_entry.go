package domain

import "time"

type QueueStatus string

const (
	QueueStatusWaiting QueueStatus = "WAITING"
	QueueStatusReady   QueueStatus = "READY"
	QueueStatusExpired QueueStatus = "EXPIRED"
)

type QueueEntry struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	EventID    uint        `gorm:"not null;index" json:"event_id"`
	ClientID   string      `gorm:"size:100;not null;index" json:"client_id"`
	QueueToken string      `gorm:"size:100;not null;uniqueIndex" json:"queue_token"`
	Status     QueueStatus `gorm:"size:30;not null" json:"status"`
	Position   int         `gorm:"not null" json:"position"`
	ReadyAt    *time.Time  `json:"ready_at"`
	ExpiredAt  *time.Time  `json:"expired_at"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

func (QueueEntry) TableName() string {
	return "queue_entries"
}
