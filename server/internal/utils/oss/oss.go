package oss

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"cpmail/internal/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSSService struct {
	client *oss.Client
	bucket *oss.Bucket
}

func NewOSSService() (*OSSService, error) {
	endpoint := fmt.Sprintf("oss-%s.aliyuncs.com", config.AppConfig.OSS.Region)

	// 创建OSS客户端
	client, err := oss.New(
		endpoint,
		config.AppConfig.OSS.AccessKeyId,
		config.AppConfig.OSS.AccessKeySecret,
	)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %v", err)
	}

	// 获取存储空间
	bucket, err := client.Bucket(config.AppConfig.OSS.Bucket)
	if err != nil {
		return nil, fmt.Errorf("获取存储空间失败: %v", err)
	}

	return &OSSService{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *OSSService) UploadFile(file *multipart.FileHeader) (string, bool, error) {

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", false, err
	}
	defer src.Close()

	// 计算文件MD5哈希值
	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", false, err
	}

	// 重置文件指针
	if _, err := src.Seek(0, 0); err != nil {
		return "", false, err
	}

	fileHash := fmt.Sprintf("%x", hash.Sum(nil))
	fileExt := filepath.Ext(file.Filename)

	// 生成文件名
	objectKey := fmt.Sprintf("uploads/%s%s", fileHash, fileExt)

	// 检查文件是否已存在
	exists, err := s.bucket.IsObjectExist(objectKey)
	if err != nil {
		return "", false, err
	}

	if exists {
		// 文件已存在，返回URL
		url := fmt.Sprintf("https://%s.%s.aliyuncs.com/%s",
			config.AppConfig.OSS.Bucket,
			fmt.Sprintf("oss-%s", config.AppConfig.OSS.Region),
			objectKey)
		return url, true, nil
	}

	// 上传文件
	err = s.bucket.PutObject(objectKey, src)
	if err != nil {
		return "", false, err
	}

	// 返回文件URL
	url := fmt.Sprintf("https://%s.%s.aliyuncs.com/%s",
		config.AppConfig.OSS.Bucket,
		fmt.Sprintf("oss-%s", config.AppConfig.OSS.Region),
		objectKey)

	return url, false, nil
}

func (s *OSSService) UploadTemplate(file *multipart.FileHeader) (string, error) {
	// 生成模板文件的object key
	objectKey := "template/todaytemplate.xlsx"

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 上传文件
	err = s.bucket.PutObject(objectKey, src)
	if err != nil {
		return "", err
	}

	// 返回文件URL
	url := fmt.Sprintf("https://%s.%s.aliyuncs.com/%s",
		config.AppConfig.OSS.Bucket,
		fmt.Sprintf("oss-%s", config.AppConfig.OSS.Region),
		objectKey)

	return url, nil
}

func (s *OSSService) GetTemplateURL() string {
	objectKey := "template/todaytemplate.xlsx"
	url := fmt.Sprintf("https://%s.%s.aliyuncs.com/%s",
		config.AppConfig.OSS.Bucket,
		fmt.Sprintf("oss-%s", config.AppConfig.OSS.Region),
		objectKey)
	return url
}

func (s *OSSService) TemplateExists() (bool, error) {
	objectKey := "template/todaytemplate.xlsx"
	return s.bucket.IsObjectExist(objectKey)
}
