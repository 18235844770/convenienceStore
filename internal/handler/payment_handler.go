package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/service"
)

// PaymentHandler 处理支付相关的回调。
type PaymentHandler struct {
	service service.PaymentService
}

// NewPaymentHandler 创建 PaymentHandler 实例。
func NewPaymentHandler(service service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

// HandleWeChatCallback 处理来自微信支付的异步通知。
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
