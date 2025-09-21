package handler

import (
	"cpmail/internal/service"
	"cpmail/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailService *service.EmailService
}

func NewEmailHandler(emailService *service.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

func (h *EmailHandler) SendEmail(c *gin.Context) {
	h.emailService.SendEmail(c)
}

func (h *EmailHandler) GetEmailStats(c *gin.Context) {
	// 获取邮件发送统计
	stats, err := h.emailService.GetDBService().GetEmailStats()
	if err != nil {
		response.Fail(c, "获取统计信息失败")
		return
	}
	
	response.Success(c, stats)
}