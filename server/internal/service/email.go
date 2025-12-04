package service

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"cpmail/internal/utils/db"
	"cpmail/internal/utils/response"
	"cpmail/internal/utils/smtp"

	"github.com/gin-gonic/gin"
)

type EmailService struct {
	// 邮件队列
	queue     chan *EmailTask
	running   bool
	mutex     sync.Mutex
	dbService *db.DBService

	// 节流与重试控制
	sendInterval time.Duration
	maxRetries   int
}

type EmailTask struct {
	Email              string                `json:"email"`
	Password           string                `json:"password"`
	Subject            string                `json:"subject"`
	ReceiverItemsArray [][]interface{}       `json:"receiverItemsArray"`
	Content            string                `json:"content"`
	ResponseChan       chan *EmailSendResult `json:"-"`
}

type EmailSendResult struct {
	SuccessList []string `json:"successList"`
}

type SendEmailRequest struct {
	Email              string          `json:"email" binding:"required"`
	Password           string          `json:"password" binding:"required"`
	Subject            string          `json:"subject" binding:"required"`
	ReceiverItemsArray [][]interface{} `json:"receiverItemsArray" binding:"required"`
	Content            string          `json:"content" binding:"required"`
}

func NewEmailService(dbService *db.DBService) *EmailService {
	service := &EmailService{
		queue:        make(chan *EmailTask, 100), // 队列大小为100
		running:      false,
		dbService:    dbService,
		sendInterval: time.Second,
		maxRetries:   2,
	}

	// 启动邮件处理goroutine
	go service.processQueue()

	return service
}

func (s *EmailService) SendEmail(c *gin.Context) {
	var req SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.INVALID_PARAMS, "参数错误")
		return
	}

	// 创建任务
	task := &EmailTask{
		Email:              req.Email,
		Password:           req.Password,
		Subject:            req.Subject,
		ReceiverItemsArray: req.ReceiverItemsArray,
		Content:            req.Content,
		ResponseChan:       make(chan *EmailSendResult, 1),
	}

	// 将任务加入队列
	select {
	case s.queue <- task:
		// 等待处理结果
		result := <-task.ResponseChan
		// 更新邮件发送统计
		s.dbService.IncrementEmailCount(len(result.SuccessList))
		response.Success(c, result)
	case <-time.After(5 * time.Second):
		response.Error(c, response.ERROR, "邮件队列已满，请稍后再试")
	}
}

func (s *EmailService) GetDBService() *db.DBService {
	return s.dbService
}

func (s *EmailService) processQueue() {
	for task := range s.queue {
		result := s.sendEmails(task)
		task.ResponseChan <- result
	}
}

func (s *EmailService) sendEmails(task *EmailTask) *EmailSendResult {
	// 选择SMTP配置
	smtpConfig := smtp.SelectSmtpConfig(task.Email)
	if smtpConfig == nil {
		return &EmailSendResult{
			SuccessList: []string{},
		}
	}

	smtpConfig.Username = task.Email
	smtpConfig.Password = task.Password

	successList := make([]string, 0, len(task.ReceiverItemsArray))
	failureQueue := task.ReceiverItemsArray

	client, err := s.connectSMTPClient(smtpConfig)
	if err != nil {
		fmt.Printf("连接SMTP服务器失败: %v\n", err)
		return &EmailSendResult{SuccessList: successList}
	}
	defer func() {
		if client != nil {
			client.Close()
		}
	}()

	sendBatch := func(client *smtp.SMTPClient, queue [][]interface{}) [][]interface{} {
		var remaining [][]interface{}
		for _, items := range queue {
			if len(items) < 3 {
				continue
			}

			receiverEmail, ok := items[0].(string)
			if !ok {
				continue
			}

			variables, ok := items[2].(map[string]interface{})
			if !ok {
				variables = make(map[string]interface{})
			}

			replacedContent := task.Content
			for varName, varValue := range variables {
				varValueStr := fmt.Sprintf("%v", varValue)
				replacedContent = strings.ReplaceAll(replacedContent, "{"+varName+"}", varValueStr)
			}

			if err := client.SendEmail([]string{receiverEmail}, task.Subject, replacedContent); err != nil {
				fmt.Printf("发送邮件给 %s 失败: %v\n", receiverEmail, err)
				remaining = append(remaining, items)
			} else {
				successList = append(successList, receiverEmail)
			}

			time.Sleep(s.sendInterval)
		}

		return remaining
	}

	failureQueue = sendBatch(client, failureQueue)
	for attempt := 1; attempt <= s.maxRetries && len(failureQueue) > 0; attempt++ {
		client.Close()
		client, err = s.connectSMTPClient(smtpConfig)
		if err != nil {
			fmt.Printf("重试连接SMTP失败(第%d次): %v\n", attempt, err)
			client = nil
			continue
		}
		failureQueue = sendBatch(client, failureQueue)
	}

	return &EmailSendResult{SuccessList: successList}
}

func (s *EmailService) connectSMTPClient(smtpConfig *smtp.SMTPConfig) (*smtp.SMTPClient, error) {
	client := smtp.NewSMTPClient(smtpConfig)
	if err := client.Connect(); err != nil {
		return nil, err
	}
	return client, nil
}
