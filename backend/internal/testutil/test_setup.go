package testutil

import (
	"context"
	"time"

	"ticketeer/backend/internal/config"
	"ticketeer/backend/internal/domain"
	"ticketeer/backend/internal/infra"
	"ticketeer/backend/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupTestApp() (*gorm.DB, *redis.Client, *gin.Engine, error) {
	gin.SetMode(gin.TestMode)

	cfg, err := config.Load()
	if err != nil {
		return nil, nil, nil, err
	}

	db, err := infra.NewPostgres(cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := db.AutoMigrate(&domain.Event{}, &domain.Seat{}, &domain.QueueEntry{}, &domain.Booking{}); err != nil {
		return nil, nil, nil, err
	}

	rdb, err := infra.NewRedis(cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	engine := router.SetupRouter(db, rdb)
	return db, rdb, engine, nil
}

func ResetState(db *gorm.DB, rdb *redis.Client) error {
	if err := db.Exec("TRUNCATE TABLE bookings, queue_entries, seats, events RESTART IDENTITY CASCADE").Error; err != nil {
		return err
	}
	return rdb.FlushDB(context.Background()).Err()
}

func SeedOpenEvent(db *gorm.DB) (*domain.Event, error) {
	now := time.Now()
	event := &domain.Event{
		Title:          "Test Concert",
		Description:    "test",
		Venue:          "Test Hall",
		EventAt:        now.Add(24 * time.Hour),
		BookingOpenAt:  now.Add(-1 * time.Hour),
		BookingCloseAt: now.Add(24 * time.Hour),
		Status:         domain.EventStatusOpen,
	}
	return event, db.Create(event).Error
}

func SeedPendingEvent(db *gorm.DB) (*domain.Event, error) {
	now := time.Now()
	event := &domain.Event{
		Title:          "Pending Concert",
		Description:    "test",
		Venue:          "Test Hall",
		EventAt:        now.Add(24 * time.Hour),
		BookingOpenAt:  now.Add(24 * time.Hour),
		BookingCloseAt: now.Add(48 * time.Hour),
		Status:         domain.EventStatusOpenPending,
	}
	return event, db.Create(event).Error
}

func SeedSeat(db *gorm.DB, eventID uint, seatNo string, status domain.SeatStatus) (*domain.Seat, error) {
	seat := &domain.Seat{
		EventID: eventID,
		SeatNo:  seatNo,
		Section: "R",
		Price:   150000,
		Status:  status,
	}
	return seat, db.Create(seat).Error
}
