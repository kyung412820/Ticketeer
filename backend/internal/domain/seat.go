package domain

import "time"

type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "AVAILABLE"
	SeatStatusHeld      SeatStatus = "HELD"
	SeatStatusBooked    SeatStatus = "BOOKED"
)

type Seat struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	EventID   uint       `gorm:"not null;index;uniqueIndex:idx_event_seat_no" json:"event_id"`
	SeatNo    string     `gorm:"size:50;not null;uniqueIndex:idx_event_seat_no" json:"seat_no"`
	Section   string     `gorm:"size:50;not null" json:"section"`
	Price     int        `gorm:"not null" json:"price"`
	Status    SeatStatus `gorm:"size:30;not null" json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Seat) TableName() string {
	return "seats"
}
