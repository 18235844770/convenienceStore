package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/model"
	"convenienceStore/internal/service"
)

// UserHandler 对外提供用户相关的 HTTP 接口。
type UserHandler struct {
	service service.UserService
}

// NewUserHandler 构建新的 UserHandler 实例。
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// WeChatLogin 使用微信授权码换取用户会话。
func (h *UserHandler) WeChatLogin(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.WeChatLogin(c.Request.Context(), req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// BindUser 将用户账号与扩展档案数据绑定。
func (h *UserHandler) BindUser(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.BindUser(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListAddresses 返回用户保存的全部收货地址。
func (h *UserHandler) ListAddresses(c *gin.Context) {
	userID := c.Query("user_id")
	addresses, err := h.service.ListAddresses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

// CreateAddress 为用户新增收货地址。
func (h *UserHandler) CreateAddress(c *gin.Context) {
	var req model.Address
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateAddress(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// UpdateAddress 更新已有的收货地址。
func (h *UserHandler) UpdateAddress(c *gin.Context) {
	var req model.Address
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = c.Param("id")

	if err := h.service.UpdateAddress(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteAddress 从用户资料中移除该地址。
func (h *UserHandler) DeleteAddress(c *gin.Context) {
	if err := h.service.DeleteAddress(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
