package handler

import (
	"errors"
	"net/http"
	"strconv"

	"ticketeer/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type SeatHandler struct {
	seatService *service.SeatService
}

func NewSeatHandler(seatService *service.SeatService) *SeatHandler {
	return &SeatHandler{
		seatService: seatService,
	}
}

func (h *SeatHandler) GetSeatsByEventID(c *gin.Context) {
	idParam := c.Param("id")
	eventID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_EVENT_ID",
				"message": "유효하지 않은 공연 ID입니다.",
			},
		})
		return
	}

	seats, err := h.seatService.GetSeatsByEventID(uint(eventID64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "좌석 목록 조회 중 오류가 발생했습니다.",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"event_id": uint(eventID64),
			"seats":    seats,
		},
	})
}

type HoldSeatRequest struct {
	EventID    uint   `json:"event_id"`
	ClientID   string `json:"client_id"`
	QueueToken string `json:"queue_token"`
}

func (h *SeatHandler) HoldSeat(c *gin.Context) {
	seatIDParam := c.Param("seatId")
	seatID64, err := strconv.ParseUint(seatIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{"code": "INVALID_SEAT_ID", "message": "유효하지 않은 좌석 ID입니다."},
		})
		return
	}

	var req HoldSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{"code": "INVALID_REQUEST", "message": "잘못된 요청입니다."},
		})
		return
	}

	result, err := h.seatService.HoldSeat(uint(seatID64), req.EventID, req.ClientID, req.QueueToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSeatNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error": gin.H{"code": "SEAT_NOT_FOUND", "message": "존재하지 않는 좌석입니다."},
			})
		case errors.Is(err, service.ErrSeatEventMismatch):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{"code": "SEAT_EVENT_MISMATCH", "message": "공연과 좌석 정보가 일치하지 않습니다."},
			})
		case errors.Is(err, service.ErrSeatAlreadyBooked):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{"code": "SEAT_ALREADY_BOOKED", "message": "이미 예매 완료된 좌석입니다."},
			})
		case errors.Is(err, service.ErrSeatAlreadyHeld):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{"code": "SEAT_ALREADY_HELD", "message": "이미 다른 사용자가 선택한 좌석입니다."},
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
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{"code": "INTERNAL_SERVER_ERROR", "message": "좌석 홀드 중 오류가 발생했습니다."},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}
