package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/service"
)

// ProductHandler exposes product-related endpoints.
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler constructs a ProductHandler instance.
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func parseStatusQuery(value string) (*bool, error) {
	if value == "" {
		return nil, nil
	}

	switch value {
	case "active":
		v := true
		return &v, nil
	case "inactive":
		v := false
		return &v, nil
	default:
		return nil, fmt.Errorf("invalid status value: %s", value)
	}
}

// ListProducts returns a list of products, optionally filtered by status.
func (h *ProductHandler) ListProducts(c *gin.Context) {
	statusFilter, err := parseStatusQuery(c.Query("status"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := h.service.ListProducts(c.Request.Context(), statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct returns a product by ID, optionally filtered by status.
func (h *ProductHandler) GetProduct(c *gin.Context) {
	statusFilter, err := parseStatusQuery(c.Query("status"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.GetProduct(c.Request.Context(), c.Param("id"), statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// ValidateInventory checks stock availability before ordering.
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
