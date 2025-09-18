package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/service"
)

// ProductHandler 提供商品相关接口。
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler 创建 ProductHandler 实例。
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// ListProducts 返回商品的分页视图。
func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.service.ListProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct 按 ID 返回商品详情。
func (h *ProductHandler) GetProduct(c *gin.Context) {
	product, err := h.service.GetProduct(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// ValidateInventory 在下单前执行库存校验。
func (h *ProductHandler) ValidateInventory(c *gin.Context) {
	var req struct {
		Quantity int `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ok, err := h.service.ValidateInventory(c.Request.Context(), c.Param("id"), req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available": ok})
}
