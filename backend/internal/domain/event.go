package domain

import "time"

type EventStatus string

const (
	EventStatusOpenPending EventStatus = "OPEN_PENDING"
	EventStatusOpen        EventStatus = "OPEN"
	EventStatusClosed      EventStatus = "CLOSED"
)

type Event struct {
	ID             uint        `gorm:"primaryKey" json:"id"`
	Title          string      `gorm:"size:255;not null" json:"title"`
	Description    string      `gorm:"type:text" json:"description"`
	Venue          string      `gorm:"size:255;not null" json:"venue"`
	EventAt        time.Time   `gorm:"not null" json:"event_at"`
	BookingOpenAt  time.Time   `gorm:"not null" json:"booking_open_at"`
	BookingCloseAt time.Time   `gorm:"not null" json:"booking_close_at"`
	Status         EventStatus `gorm:"size:30;not null" json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

func (Event) TableName() string {
	return "events"
}
