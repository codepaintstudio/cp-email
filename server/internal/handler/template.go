package handler

import (
	"cpmail/internal/service"
	"cpmail/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	templateService *service.TemplateService
}

func NewTemplateHandler(templateService *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
	}
}

func (h *TemplateHandler) DownloadTemplate(c *gin.Context) {
	// 检查模板文件是否存在
	if !h.templateService.TemplateExists() {
		response.Error(c, response.ERROR, "模板文件不存在，请检查路径或文件是否已上传")
		return
	}

	// 重定向到OSS文件URL
	c.Redirect(302, h.templateService.GetTemplateURL())
}

func (h *TemplateHandler) UploadTemplate(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, response.INVALID_PARAMS, "请选择要上传的文件")
		return
	}

	// 保存文件到OSS
	_, err = h.templateService.UploadTemplate(file)
	if err != nil {
		response.Fail(c, "文件保存失败")
		return
	}

	response.Success(c, gin.H{
		"msg": "模板上传成功",
	})
}