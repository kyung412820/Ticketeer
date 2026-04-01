package router

import (
	"time"

	"ticketeer/backend/internal/handler"
	"ticketeer/backend/internal/repository"
	"ticketeer/backend/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	healthHandler := handler.NewHealthHandler()

	eventRepository := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepository)
	eventHandler := handler.NewEventHandler(eventService)

	queueRepository := repository.NewQueueRepository(db)
	queueService := service.NewQueueService(eventRepository, queueRepository)
	queueHandler := handler.NewQueueHandler(queueService)

	seatRepository := repository.NewSeatRepository(db)
	seatService := service.NewSeatService(seatRepository, queueRepository, rdb)
	seatHandler := handler.NewSeatHandler(seatService)

	bookingRepository := repository.NewBookingRepository(db)
	bookingService := service.NewBookingService(db, seatRepository, queueRepository, bookingRepository, rdb)
	bookingHandler := handler.NewBookingHandler(bookingService)

	api := r.Group("/api")
	{
		api.GET("/health", healthHandler.Check)
		api.GET("/events", eventHandler.GetEvents)
		api.GET("/events/:id", eventHandler.GetEvent)
		api.GET("/events/:id/seats", seatHandler.GetSeatsByEventID)
		api.POST("/queue/enter", queueHandler.Enter)
		api.GET("/queue/status/:token", queueHandler.GetStatus)
		api.POST("/seats/:seatId/hold", seatHandler.HoldSeat)
		api.POST("/bookings", bookingHandler.Create)
	}

	return r
}
