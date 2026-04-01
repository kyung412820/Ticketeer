package handler

import (
	"errors"
	"net/http"

	"ticketeer/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

type CreateBookingRequest struct {
	EventID    uint   `json:"event_id"`
	SeatID     uint   `json:"seat_id"`
	ClientID   string `json:"client_id"`
	QueueToken string `json:"queue_token"`
}

func (h *BookingHandler) Create(c *gin.Context) {
	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{"code": "INVALID_REQUEST", "message": "잘못된 요청입니다."},
		})
		return
	}

	result, err := h.bookingService.CreateBooking(req.EventID, req.SeatID, req.ClientID, req.QueueToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSeatNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error": gin.H{"code": "SEAT_NOT_FOUND", "message": "존재하지 않는 좌석입니다."},
			})
		case errors.Is(err, service.ErrSeatAlreadyBooked):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{"code": "SEAT_ALREADY_BOOKED", "message": "이미 예매 완료된 좌석입니다."},
			})
		case errors.Is(err, service.ErrQueueTokenNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error": gin.H{"code": "QUEUE_TOKEN_NOT_FOUND", "message": "존재하지 않는 대기열 토큰입니다."},
			})
		case errors.Is(err, service.ErrInvalidQueueToken):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{"code": "INVALID_QUEUE_TOKEN", "message": "유효하지 않은 대기열 토큰입니다."},
			})
		case errors.Is(err, service.ErrQueueExpired):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{"code": "QUEUE_EXPIRED", "message": "대기열 토큰이 만료되었습니다."},
			})
		case errors.Is(err, service.ErrHoldNotFound):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{"code": "HOLD_NOT_FOUND", "message": "좌석 홀드 정보가 없습니다."},
			})
		case errors.Is(err, service.ErrHoldExpired):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{"code": "HOLD_EXPIRED", "message": "좌석 홀드 시간이 만료되었습니다."},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{"code": "INTERNAL_SERVER_ERROR", "message": "예매 확정 중 오류가 발생했습니다."},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}
