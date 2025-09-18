package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/model"
	"convenienceStore/internal/service"
)

// OrderHandler 提供订单生命周期相关的接口。
type OrderHandler struct {
	service service.OrderService
}

// NewOrderHandler 实例化 OrderHandler。
func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrder 生成新的订单草稿。
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req model.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.service.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder 返回订单详情。
func (h *OrderHandler) GetOrder(c *gin.Context) {
	order, err := h.service.GetOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// PayOrder 发起订单支付。
func (h *OrderHandler) PayOrder(c *gin.Context) {
	paymentInfo, err := h.service.PayOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentInfo)
}

// CancelOrder 在发货前取消订单。
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	if err := h.service.CancelOrder(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ShipOrder 将订单状态更新为已发货。
func (h *OrderHandler) ShipOrder(c *gin.Context) {
	if err := h.service.ShipOrder(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CompleteOrder 将订单标记为已完成。
func (h *OrderHandler) CompleteOrder(c *gin.Context) {
	if err := h.service.CompleteOrder(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
