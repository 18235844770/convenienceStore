package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/service"
)

// AdminProductHandler exposes management endpoints for products.
type AdminProductHandler struct {
	service service.AdminProductService
}

// NewAdminProductHandler constructs an AdminProductHandler instance.
func NewAdminProductHandler(service service.AdminProductService) *AdminProductHandler {
	return &AdminProductHandler{service: service}
}

type adminProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required"`
	Stock       int      `json:"stock" binding:"required"`
	Tags        []string `json:"tags"`
	Images      []string `json:"images"`
	IsActive    *bool    `json:"is_active"`
}

// ListProducts returns products for the management console.
func (h *AdminProductHandler) ListProducts(c *gin.Context) {
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

// GetProduct returns a single product for the management console.
func (h *AdminProductHandler) GetProduct(c *gin.Context) {
	product, err := h.service.GetProduct(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct creates a new product record.
func (h *AdminProductHandler) CreateProduct(c *gin.Context) {
	var req adminProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := service.AdminProductPayload{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Tags:        req.Tags,
		Images:      req.Images,
		IsActive:    req.IsActive,
	}

	product, err := h.service.CreateProduct(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct updates an existing product record.
func (h *AdminProductHandler) UpdateProduct(c *gin.Context) {
	var req adminProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := service.AdminProductPayload{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Tags:        req.Tags,
		Images:      req.Images,
		IsActive:    req.IsActive,
	}

	product, err := h.service.UpdateProduct(c.Request.Context(), c.Param("id"), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct removes a product record.
func (h *AdminProductHandler) DeleteProduct(c *gin.Context) {
	if err := h.service.DeleteProduct(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SetProductStatus toggles the availability of a product.
func (h *AdminProductHandler) SetProductStatus(c *gin.Context) {
	var req struct {
		IsActive bool `json:"is_active" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SetProductStatus(c.Request.Context(), c.Param("id"), req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.GetProduct(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}
