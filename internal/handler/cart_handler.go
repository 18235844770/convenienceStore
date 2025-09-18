package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/model"
	"convenienceStore/internal/service"
)

// CartHandler 负责购物车相关的接口处理。
type CartHandler struct {
	service service.CartService
}

// NewCartHandler 创建 CartHandler 实例。
func NewCartHandler(service service.CartService) *CartHandler {
	return &CartHandler{service: service}
}

// ListItems 返回用户的购物车内容。
func (h *CartHandler) ListItems(c *gin.Context) {
	userID := c.Query("user_id")
	items, err := h.service.ListItems(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// AddItem 向购物车新增商品。
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

// UpdateItem 调整购物车条目的数量或选中状态。
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

// RemoveItem 从购物车移除商品。
func (h *CartHandler) RemoveItem(c *gin.Context) {
	if err := h.service.RemoveItem(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
