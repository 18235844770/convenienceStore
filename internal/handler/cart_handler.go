package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/model"
	"convenienceStore/internal/service"
)

// CartHandler manages shopping cart endpoints.
type CartHandler struct {
	service service.CartService
}

// NewCartHandler creates a CartHandler instance.
func NewCartHandler(service service.CartService) *CartHandler {
	return &CartHandler{service: service}
}

// ListItems returns the cart content for a user.
func (h *CartHandler) ListItems(c *gin.Context) {
	userID := c.Query("user_id")
	items, err := h.service.ListItems(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// AddItem adds a new item to the cart.
func (h *CartHandler) AddItem(c *gin.Context) {
	var item model.CartItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddItem(c.Request.Context(), &item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// UpdateItem adjusts quantity or selection status of a cart entry.
func (h *CartHandler) UpdateItem(c *gin.Context) {
	var item model.CartItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item.ID = c.Param("id")

	if err := h.service.UpdateItem(c.Request.Context(), &item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveItem removes an item from the cart.
func (h *CartHandler) RemoveItem(c *gin.Context) {
	if err := h.service.RemoveItem(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
