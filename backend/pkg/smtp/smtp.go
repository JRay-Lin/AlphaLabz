package smtp

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type EmailService struct {
	config EmailConfig
}

func NewEmailService(config EmailConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

func (s *EmailService) SendEmail(to []string, subject string, templateName string, data interface{}) error {
	// Load template
	tmpl, err := template.ParseFiles(fmt.Sprintf("pkg/smtp/template/%s.html", templateName))
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// Compose email
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	headers := fmt.Sprintf("To: %s\r\nSubject: %s\r\n%s\r\n", to[0], subject, mime)
	message := []byte(headers + body.String())

	// Send email
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	if err := smtp.SendMail(addr, auth, s.config.From, to, message); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
