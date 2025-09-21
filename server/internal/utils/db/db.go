package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"cpmail/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBService struct {
	db *gorm.DB
}

type EmailStat struct {
	ID          uint      `gorm:"primaryKey"`
	TotalCount  int       `gorm:"default:0"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func NewDBService() *DBService {
	// 确保数据库目录存在
	dbDir := filepath.Dir(config.AppConfig.Database.Path)
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		log.Fatalf("创建数据库目录失败: %v", err)
	}

	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(config.AppConfig.Database.Path), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移数据库结构
	if err := db.AutoMigrate(&EmailStat{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 确保有一条统计记录
	var stat EmailStat
	if err := db.FirstOrCreate(&stat, EmailStat{ID: 1}).Error; err != nil {
		log.Fatalf("初始化统计记录失败: %v", err)
	}

	return &DBService{
		db: db,
	}
}

func (s *DBService) IncrementEmailCount(count int) error {
	if count <= 0 {
		return nil
	}

	// 更新邮件发送总数
	if err := s.db.Model(&EmailStat{}).Where("id = ?", 1).UpdateColumn("total_count", gorm.Expr("total_count + ?", count)).Error; err != nil {
		return fmt.Errorf("更新邮件统计失败: %v", err)
	}

	return nil
}

func (s *DBService) GetEmailStats() (*EmailStat, error) {
	var stat EmailStat
	if err := s.db.First(&stat, 1).Error; err != nil {
		return nil, fmt.Errorf("获取邮件统计失败: %v", err)
	}
	return &stat, nil
}