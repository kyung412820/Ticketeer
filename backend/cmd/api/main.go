package main

import (
	"log"

	"ticketeer/backend/internal/config"
	"ticketeer/backend/internal/domain"
	"ticketeer/backend/internal/infra"
	"ticketeer/backend/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := infra.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}

	if err := db.AutoMigrate(&domain.Event{}, &domain.Seat{}, &domain.QueueEntry{}, &domain.Booking{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	rdb, err := infra.NewRedis(cfg)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	r := router.SetupRouter(db, rdb)

	log.Printf("server started on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
