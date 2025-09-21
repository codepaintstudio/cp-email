package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"cpmail/internal/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	jwtSecret  string
	jwtExpires time.Duration
}

func NewAuthService(secret string, expires string) *AuthService {
	duration, _ := time.ParseDuration(expires)
	return &AuthService{
		jwtSecret:  secret,
		jwtExpires: duration,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *AuthService) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.INVALID_PARAMS, "参数错误")
		return
	}

	// 验证邮箱和密码
	if err := s.verifyEmailCredentials(req.Email, req.Password); err != nil {
		response.Error(c, response.ACCOUNT_OR_PASSWORD_ERROR, "请检查账号密码是否正确！")
		return
	}

	// 生成 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": req.Email,
		"exp":   time.Now().Add(s.jwtExpires).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		response.Fail(c, "生成token失败")
		return
	}

	response.Success(c, gin.H{
		"msg":   "账号密码验证成功",
		"token": "Bearer " + tokenString,
	})
}

func (s *AuthService) verifyEmailCredentials(email, password string) error {
	// 根据邮箱域名选择 SMTP 配置
	smtpConfig := s.selectSmtpConfig(email)
	if smtpConfig == nil {
		return fmt.Errorf("暂不支持此类邮箱")
	}

	fmt.Printf("尝试连接 SMTP 服务器: %s:%d\n", smtpConfig.host, smtpConfig.port)
	
	// 创建 SMTP 客户端
	auth := smtp.PlainAuth("", email, password, smtpConfig.host)
	
	// 对于465端口，使用SSL连接
	if smtpConfig.port == 465 {
		fmt.Println("使用 SSL 连接")
		tlsConfig := &tls.Config{
			ServerName: smtpConfig.host,
		}
		conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", smtpConfig.host, smtpConfig.port), tlsConfig)
		if err != nil {
			fmt.Printf("SSL 连接失败: %v\n", err)
			return err
		}
		defer conn.Close()
		
		client, err := smtp.NewClient(conn, smtpConfig.host)
		if err != nil {
			fmt.Printf("创建 SMTP 客户端失败: %v\n", err)
			return err
		}
		defer client.Close()
		
		fmt.Println("尝试认证")
		// 验证认证信息
		if err = client.Auth(auth); err != nil {
			fmt.Printf("认证失败: %v\n", err)
			return err
		}
		fmt.Println("认证成功")
	} else {
		fmt.Println("使用普通连接")
		// 对于其他端口，使用普通连接
		client, err := smtp.Dial(fmt.Sprintf("%s:%d", smtpConfig.host, smtpConfig.port))
		if err != nil {
			fmt.Printf("连接失败: %v\n", err)
			return err
		}
		defer client.Close()
		
		// 如果支持StartTLS，则使用它
		if smtpConfig.port == 587 {
			if ok, _ := client.Extension("STARTTLS"); ok {
				fmt.Println("使用 STARTTLS")
				tlsConfig := &tls.Config{
					ServerName: smtpConfig.host,
				}
				if err = client.StartTLS(tlsConfig); err != nil {
					fmt.Printf("STARTTLS 失败: %v\n", err)
					return err
				}
			}
		}
		
		fmt.Println("尝试认证")
		// 验证认证信息
		if err = client.Auth(auth); err != nil {
			fmt.Printf("认证失败: %v\n", err)
			return err
		}
		fmt.Println("认证成功")
	}

	return nil
}

type smtpConfig struct {
	host string
	port int
}

func (s *AuthService) selectSmtpConfig(email string) *smtpConfig {
	realmName := strings.Split(email, "@")[1]
	switch realmName {
	case "163.com":
		return &smtpConfig{
			host: "smtp.163.com",
			port: 465,
		}
	case "qq.com":
		return &smtpConfig{
			host: "smtp.qq.com",
			port: 465,
		}
	default:
		return nil
	}
}