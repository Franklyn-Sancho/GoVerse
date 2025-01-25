package email

import (
	"fmt"
	"net/smtp"
)

type EmailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailService interface {
	SendEmail(to, subject, body string) error
	SendConfirmationEmail(email, username, token string) error
}

type emailService struct {
	queueService EmailQueueService
}

func NewEmailService(queueService EmailQueueService) EmailService {
	return &emailService{queueService: queueService}
}

func (s *emailService) SendEmail(to, subject, body string) error {
	from := "test@test.com"
	smtpHost := "mailhog"
	smtpPort := "1025"

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))
	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, msg)
}

func (s *emailService) SendConfirmationEmail(email, username, token string) error {
	msg := EmailMessage{
		To:      email,
		Subject: "Confirmação de Registro",
		Body:    fmt.Sprintf("Olá %s,\n\nConfirme seu email:\nhttp://localhost:8080/confirm-email?token=%s", username, token),
	}

	return s.queueService.PublishEmail(msg)
}
