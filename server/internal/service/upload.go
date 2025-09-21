package service

import (
	"fmt"
	"log"
	"mime/multipart"

	"cpmail/internal/utils/oss"
)

type UploadService struct {
	ossService *oss.OSSService
}

func NewUploadService() *UploadService {
	log.Printf("[Upload Service] 开始初始化上传服务")

	// 初始化OSS服务
	ossService, err := oss.NewOSSService()
	if err != nil {
		log.Printf("[Upload Service] 初始化OSS服务失败: %v", err)
		// 返回一个空的服务实例，避免nil指针
		return &UploadService{
			ossService: nil,
		}
	}

	log.Printf("[Upload Service] OSS服务初始化成功")
	return &UploadService{
		ossService: ossService,
	}
}

func (s *UploadService) UploadFile(file *multipart.FileHeader) (string, bool, error) {
	log.Printf("[Upload Service] 开始上传文件: %s", file.Filename)

	// 检查OSS服务是否初始化成功
	if s.ossService == nil {
		log.Printf("[Upload Service] OSS服务未初始化")
		return "", false, fmt.Errorf("OSS服务未初始化")
	}

	log.Printf("[Upload Service] OSS服务已初始化，调用OSS上传")

	// 使用OSS上传文件
	url, isExisting, err := s.ossService.UploadFile(file)
	if err != nil {
		log.Printf("[Upload Service] OSS上传失败: %v", err)
		return "", false, err
	}

	log.Printf("[Upload Service] OSS上传成功: URL=%s, isExisting=%v", url, isExisting)
	return url, isExisting, nil
}