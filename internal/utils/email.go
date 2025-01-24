package utils

import (
	"fmt"
	"net/smtp"
)

func SendEmail(to string, subject string, body string) error {
	from := "your_email@example.com"

	smtpHost := "mailhog"
	smtpPort := "1025"

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, msg)
}

func SendConfirmationEmail(email, username, token string) error {
	subject := "Confirmação de Registro"
	body := fmt.Sprintf(
		"Olá %s, \n\nPor favor, clique no link abaixo para confirmar seu e-mail:\nhttp://localhost:8080/confirm-email?token=%s",
		username, token,
	)
	return SendEmail(email, subject, body)
}
