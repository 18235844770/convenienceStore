package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/model"
	"convenienceStore/internal/service"
)

// DeliveryHandler 负责处理配送相关的接口。
type DeliveryHandler struct {
	service service.DeliveryService
}

// NewDeliveryHandler 构建 DeliveryHandler 实例。
func NewDeliveryHandler(service service.DeliveryService) *DeliveryHandler {
	return &DeliveryHandler{service: service}
}

// BindAddress 将地址绑定到后续配送。
func (h *DeliveryHandler) BindAddress(c *gin.Context) {
	var address model.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.BindAddress(c.Request.Context(), &address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ShipOrder 触发订单的物流履约。
func (h *DeliveryHandler) ShipOrder(c *gin.Context) {
	var req struct {
		OrderID string `json:"order_id" binding:"required"`
		Carrier string `json:"carrier"`
		TrackNo string `json:"track_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ShipOrder(c.Request.Context(), req.OrderID, req.Carrier, req.TrackNo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
