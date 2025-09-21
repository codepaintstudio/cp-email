package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
	OSS      OSSConfig
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	SecretKey string
	ExpiresIn string
}

type DatabaseConfig struct {
	Path string
}

type OSSConfig struct {
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	Region          string
}

var AppConfig *Config

func LoadConfig() {
	// 获取当前文件的目录
	_, filename, _, _ := runtime.Caller(0)
	configPath := filepath.Join(filepath.Dir(filename), "../../configs")

	// 设置默认值
	viper.SetDefault("server.port", "4443")
	viper.SetDefault("jwt.secretKey", "mima")
	viper.SetDefault("jwt.expiresIn", "24h")
	viper.SetDefault("database.path", "./data/cpmail.db")

	// 加载配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("警告: 未找到配置文件: %v", err)
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port: viper.GetString("server.port"),
		},
		JWT: JWTConfig{
			SecretKey: viper.GetString("jwt.secretKey"),
			ExpiresIn: viper.GetString("jwt.expiresIn"),
		},
		Database: DatabaseConfig{
			Path: viper.GetString("database.path"),
		},
		OSS: OSSConfig{
			AccessKeyId:     viper.GetString("aliyun.accessKeyId"),
			AccessKeySecret: viper.GetString("aliyun.accessKeySecret"),
			Bucket:          viper.GetString("aliyun.bucket"),
			Region:          viper.GetString("aliyun.region"),
		},
	}
	
	// 输出配置信息用于调试
	log.Printf("OSS配置: AccessKeyId=%s, Bucket=%s, Region=%s", 
		AppConfig.OSS.AccessKeyId, 
		AppConfig.OSS.Bucket, 
		AppConfig.OSS.Region)
}