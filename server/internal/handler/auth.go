package handler

import (
	"cpmail/internal/service"
	"cpmail/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	h.authService.Login(c)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT 本身是无状态的，服务端不需要特殊的登出逻辑
	// 前端删除 token 即可实现登出
	response.Success(c, gin.H{
		"msg": "登出成功",
	})
}