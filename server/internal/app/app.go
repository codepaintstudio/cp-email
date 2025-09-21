package app

import (
	"log"
	"os"

	"cpmail/internal/config"
	"cpmail/internal/handler"
	"cpmail/internal/middleware"
	"cpmail/internal/service"
	"cpmail/internal/utils/db"
	"github.com/gin-gonic/gin"
)

func Run() {
	// 加载配置
	config.LoadConfig()

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化数据库服务
	dbService := db.NewDBService()

	// 初始化服务
	authService := service.NewAuthService(config.AppConfig.JWT.SecretKey, config.AppConfig.JWT.ExpiresIn)
	emailService := service.NewEmailService(dbService)
	templateService := service.NewTemplateService()
	uploadService := service.NewUploadService()

	// 初始化处理程序
	authHandler := handler.NewAuthHandler(authService)
	emailHandler := handler.NewEmailHandler(emailService)
	templateHandler := handler.NewTemplateHandler(templateService)
	uploadHandler := handler.NewUploadHandler(uploadService)

	// 初始化 JWT 中间件
	middleware.InitJWT(config.AppConfig.JWT.SecretKey)

	// 创建 Gin 引擎
	r := gin.Default()

	// 设置 CORS 中间件
	r.Use(corsMiddleware())

	// 设置静态文件服务
	r.Static("/public", "./public")

	// 注册路由
	registerRoutes(r, authHandler, emailHandler, templateHandler, uploadHandler)

	// 启动服务器
	port := config.AppConfig.Server.Port
	log.Printf("服务器启动在端口 %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
		os.Exit(1)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func registerRoutes(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	emailHandler *handler.EmailHandler,
	templateHandler *handler.TemplateHandler,
	uploadHandler *handler.UploadHandler,
) {
	// 登录路由
	loginGroup := r.Group("/api/login")
	{
		loginGroup.POST("/verify", authHandler.Login)
		loginGroup.POST("/logout", authHandler.Logout)
	}

	// 受保护的路由组
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// 邮件路由
		emailGroup := protected.Group("/api/email")
		{
			emailGroup.POST("/sendemail", emailHandler.SendEmail)
			emailGroup.GET("/stats", emailHandler.GetEmailStats)
		}

		// 模板路由
		templateGroup := protected.Group("/api/template")
		{
			templateGroup.GET("/xlsx", templateHandler.DownloadTemplate)
			templateGroup.POST("/upload", templateHandler.UploadTemplate)
		}

		// 上传路由
		uploadGroup := protected.Group("/api/uploads")
		{
			uploadGroup.POST("/file", uploadHandler.UploadFile)
		}
	}
}