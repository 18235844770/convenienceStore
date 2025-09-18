package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/service"
)

// ProductHandler exposes product related endpoints.
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler creates a ProductHandler instance.
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// ListProducts returns a paginated view of products.
func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.service.ListProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct returns product details by ID.
func (h *ProductHandler) GetProduct(c *gin.Context) {
	product, err := h.service.GetProduct(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// ValidateInventory performs a stock check before order placement.
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
