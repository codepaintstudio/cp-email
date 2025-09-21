package app

import (
	"log"
	"os"
	"path/filepath"

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

	// 获取当前工作目录，用于调试
	workDir, _ := os.Getwd()
	log.Printf("当前工作目录: %s", workDir)

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
	r.Static("/assets", "./dist/assets")
	r.StaticFile("/favicon.ico", "./dist/favicon.ico")

	// 注册路由
	registerRoutes(r, authHandler, emailHandler, templateHandler, uploadHandler)

	// 前端路由处理 - 放在最后，作为fallback
	r.NoRoute(func(c *gin.Context) {
		log.Printf("NoRoute处理路径: %s", c.Request.URL.Path)

		// 如果是API请求，返回404
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			log.Printf("API请求未找到: %s", c.Request.URL.Path)
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}

		// 检查文件是否存在，使用绝对路径
		workDir, _ := os.Getwd()
		indexPath := filepath.Join(workDir, "dist", "index.html")
		log.Printf("尝试访问文件: %s", indexPath)

		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			log.Printf("前端文件不存在: %s", indexPath)
			// 尝试备用路径
			indexPath = "./dist/index.html"
			if _, err := os.Stat(indexPath); os.IsNotExist(err) {
				log.Printf("备用路径也不存在: %s", indexPath)
				c.String(404, "前端文件未找到")
				return
			}
		}

		log.Printf("返回前端页面: %s", indexPath)
		// 返回前端index.html，让前端路由处理
		c.File(indexPath)
	})

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