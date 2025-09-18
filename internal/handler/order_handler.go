package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/model"
	"convenienceStore/internal/service"
)

// OrderHandler exposes order lifecycle endpoints.
type OrderHandler struct {
	service service.OrderService
}

// NewOrderHandler instantiates an OrderHandler.
func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrder generates a new order draft.
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

// GetOrder returns order details.
func (h *OrderHandler) GetOrder(c *gin.Context) {
	order, err := h.service.GetOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// PayOrder initiates payment for the order.
func (h *OrderHandler) PayOrder(c *gin.Context) {
	paymentInfo, err := h.service.PayOrder(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentInfo)
}

// CancelOrder cancels the order before shipment.
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	if err := h.service.CancelOrder(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ShipOrder transitions the order into shipped state.
func (h *OrderHandler) ShipOrder(c *gin.Context) {
	if err := h.service.ShipOrder(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CompleteOrder marks the order as completed.
func (h *OrderHandler) CompleteOrder(c *gin.Context) {
	if err := h.service.CompleteOrder(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
