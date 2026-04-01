package domain

import "time"

type BookingStatus string

const (
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
)

type Booking struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	EventID     uint          `gorm:"not null;index" json:"event_id"`
	SeatID      uint          `gorm:"not null;uniqueIndex" json:"seat_id"`
	ClientID    string        `gorm:"size:100;not null;index" json:"client_id"`
	BookingCode string        `gorm:"size:100;not null;uniqueIndex" json:"booking_code"`
	Status      BookingStatus `gorm:"size:30;not null" json:"status"`
	BookedAt    time.Time     `gorm:"not null" json:"booked_at"`
	CreatedAt   time.Time     `json:"created_at"`
}

func (Booking) TableName() string {
	return "bookings"
}
