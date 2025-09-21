package handler

import (
	"cpmail/internal/service"
	"cpmail/internal/utils/response"
	"log"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService *service.UploadService
}

func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	log.Printf("[Upload Handler] 开始处理文件上传请求")

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("[Upload Handler] 获取上传文件失败: %v", err)
		response.Error(c, response.INVALID_PARAMS, "请选择要上传的文件")
		return
	}

	log.Printf("[Upload Handler] 成功获取上传文件: %s, 大小: %d bytes", file.Filename, file.Size)

	// 保存文件
	fileURL, isExisting, err := h.uploadService.UploadFile(file)
	if err != nil {
		log.Printf("[Upload Handler] 文件上传失败: %v", err)
		response.Fail(c, "文件上传失败")
		return
	}

	log.Printf("[Upload Handler] 文件上传成功: URL=%s, isExisting=%v", fileURL, isExisting)

	response.Success(c, gin.H{
		"file":       fileURL,
		"isExisting": isExisting,
	})
}