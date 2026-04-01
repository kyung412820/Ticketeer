package handler

import (
	"errors"
	"net/http"

	"ticketeer/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type QueueHandler struct {
	queueService *service.QueueService
}

func NewQueueHandler(queueService *service.QueueService) *QueueHandler {
	return &QueueHandler{queueService: queueService}
}

type EnterQueueRequest struct {
	EventID  uint   `json:"event_id"`
	ClientID string `json:"client_id"`
}

func (h *QueueHandler) Enter(c *gin.Context) {
	var req EnterQueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "잘못된 요청입니다.",
			},
		})
		return
	}

	result, err := h.queueService.Enter(req.EventID, req.ClientID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEventNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error": gin.H{"code": "EVENT_NOT_FOUND", "message": "존재하지 않는 공연입니다."},
			})
		case errors.Is(err, service.ErrBookingNotOpen):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{"code": "BOOKING_NOT_OPEN", "message": "아직 예매가 시작되지 않았습니다."},
			})
		case errors.Is(err, service.ErrEventClosed):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{"code": "EVENT_CLOSED", "message": "예매가 종료된 공연입니다."},
			})
		case errors.Is(err, service.ErrAlreadyInQueue):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{"code": "ALREADY_IN_QUEUE", "message": "이미 대기열에 진입한 사용자입니다."},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{"code": "INTERNAL_SERVER_ERROR", "message": "대기열 진입 중 오류가 발생했습니다."},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (h *QueueHandler) GetStatus(c *gin.Context) {
	token := c.Param("token")
	result, err := h.queueService.GetStatus(token)
	if err != nil {
		if errors.Is(err, service.ErrQueueTokenNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error": gin.H{"code": "QUEUE_TOKEN_NOT_FOUND", "message": "존재하지 않는 대기열 토큰입니다."},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{"code": "INTERNAL_SERVER_ERROR", "message": "대기열 상태 조회 중 오류가 발생했습니다."},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}
