package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/service"
)

// UploadHandler exposes endpoints for file uploads.
type UploadHandler struct {
	service service.UploadService
}

// NewUploadHandler constructs an UploadHandler instance.
func NewUploadHandler(service service.UploadService) *UploadHandler {
	return &UploadHandler{service: service}
}

// UploadFile handles multipart file uploads and returns the stored path.
func (h *UploadHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file field is required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	path, err := h.service.SaveFile(c.Request.Context(), file.Filename, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"path": path})
}
