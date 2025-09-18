package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/service"
)

// PaymentHandler deals with payment related callbacks.
type PaymentHandler struct {
	service service.PaymentService
}

// NewPaymentHandler creates a PaymentHandler instance.
func NewPaymentHandler(service service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

// HandleWeChatCallback processes async payment notifications from WeChat Pay.
func (h *PaymentHandler) HandleWeChatCallback(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.HandleWeChatCallback(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, "success")
}
