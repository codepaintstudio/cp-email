package smtp

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type SMTPClient struct {
	config *SMTPConfig
	client *smtp.Client
}

func NewSMTPClient(config *SMTPConfig) *SMTPClient {
	return &SMTPClient{
		config: config,
	}
}

func (c *SMTPClient) Connect() error {
	var err error
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	if c.config.Port == 465 {
		// SSL connection
		conn, err := tls.Dial("tcp", addr, &tls.Config{
			ServerName: c.config.Host,
		})
		if err != nil {
			return err
		}
		c.client, err = smtp.NewClient(conn, c.config.Host)
	} else {
		// Non-SSL connection
		c.client, err = smtp.Dial(addr)
		if err != nil {
			return err
		}

		// Start TLS if available
		if c.config.Port == 587 {
			if ok, _ := c.client.Extension("STARTTLS"); ok {
				tlsConfig := &tls.Config{
					ServerName: c.config.Host,
				}
				if err = c.client.StartTLS(tlsConfig); err != nil {
					return err
				}
			}
		}
	}

	if err != nil {
		return err
	}

	// Authenticate
	auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)
	if err = c.client.Auth(auth); err != nil {
		return err
	}

	return nil
}

func (c *SMTPClient) SendEmail(to []string, subject, body string) error {
	if c.client == nil {
		return fmt.Errorf("SMTP client not connected")
	}

	from := mail.Address{
		Address: c.config.Username,
	}

	// Set sender
	if err := c.client.Mail(from.Address); err != nil {
		return err
	}

	// Set recipients
	for _, recipient := range to {
		if err := c.client.Rcpt(recipient); err != nil {
			return err
		}
	}

	// Create message
	writer, err := c.client.Data()
	if err != nil {
		return err
	}

	message := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		from.String(),
		strings.Join(to, ", "),
		subject,
		body,
	)

	_, err = writer.Write([]byte(message))
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *SMTPClient) Close() error {
	if c.client != nil {
		return c.client.Quit()
	}
	return nil
}

func SelectSmtpConfig(email string) *SMTPConfig {
	realmName := strings.Split(email, "@")[1]
	switch realmName {
	case "163.com":
		return &SMTPConfig{
			Host: "smtp.163.com",
			Port: 465,
		}
	case "qq.com":
		return &SMTPConfig{
			Host: "smtp.qq.com",
			Port: 465,
		}
	default:
		return nil
	}
}