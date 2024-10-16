package utils

import (
	"fmt"
	"net/smtp"
)

func SendEmail(to string, subject string, body string) error {
	from := "your_email@example.com" // pode ser qualquer endereço de e-mail

	// Configura o servidor SMTP do MailHog
	smtpHost := "mailhog" // MailHog normalmente está rodando no localhost
	smtpPort := "1025"    // porta do MailHog

	// Cria a mensagem
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	// Envia o e-mail sem autenticação
	return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, msg)
}

// SendConfirmationEmail envia um e-mail de confirmação para o usuário
func SendConfirmationEmail(email, username, token string) error {
	subject := "Confirmação de Registro"
	body := fmt.Sprintf(
		"Olá %s, \n\nPor favor, clique no link abaixo para confirmar seu e-mail:\nhttp://localhost:8080/confirm-email?token=%s",
		username, token,
	)
	return SendEmail(email, subject, body)
}
