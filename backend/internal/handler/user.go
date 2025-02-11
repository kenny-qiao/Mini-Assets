package handler

import (
	"mini-assets/internal/model"
	"mini-assets/internal/service"
	"mini-assets/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler 处理用户相关的HTTP请求
type UserHandler struct {
	UserService service.UserService
}

// NewUserHandler 创建UserHandler实例
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// RegisterUser 用户注册
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// 这里应该添加密码加密逻辑
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	user.Password = hashedPassword

	if err := h.UserService.RegisterUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// LoginUser 用户登录
func (h *UserHandler) LoginUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	token, err := h.UserService.LoginUser(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// SetupUserRoutes 设置用户相关的路由
func SetupUserRoutes(r *gin.Engine, userService service.UserService) {
	h := NewUserHandler(userService)
	r.POST("/register", h.RegisterUser)
	r.POST("/login", h.LoginUser)
}
