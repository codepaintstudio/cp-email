package service

import (
	"fmt"
	"mime/multipart"

	"cpmail/internal/utils/oss"
)

type TemplateService struct {
	ossService *oss.OSSService
}

func NewTemplateService() *TemplateService {
	// 初始化OSS服务
	ossService, err := oss.NewOSSService()
	if err != nil {
		fmt.Printf("初始化OSS服务失败: %v\n", err)
		// 返回一个空的服务实例，避免nil指针
		return &TemplateService{
			ossService: nil,
		}
	}
	
	return &TemplateService{
		ossService: ossService,
	}
}

func (s *TemplateService) GetTemplateURL() string {
	// 检查OSS服务是否初始化成功
	if s.ossService == nil {
		return ""
	}
	
	return s.ossService.GetTemplateURL()
}

func (s *TemplateService) TemplateExists() bool {
	// 检查OSS服务是否初始化成功
	if s.ossService == nil {
		fmt.Printf("OSS服务未初始化\n")
		return false
	}
	
	exists, err := s.ossService.TemplateExists()
	if err != nil {
		fmt.Printf("检查模板文件是否存在时出错: %v\n", err)
		return false
	}
	return exists
}

func (s *TemplateService) UploadTemplate(file *multipart.FileHeader) (string, error) {
	// 检查OSS服务是否初始化成功
	if s.ossService == nil {
		return "", fmt.Errorf("OSS服务未初始化")
	}
	
	return s.ossService.UploadTemplate(file)
}