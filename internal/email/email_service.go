package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strconv"
)

// EmailService lida com o envio de e-mails.
type EmailService struct {
	Host       string
	Port       int
	Username   string
	Password   string
	FromName   string
	FromEmail  string
	auth       smtp.Auth
}

// NewEmailService cria uma nova instância do EmailService.
func NewEmailService() (*EmailService, error) {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("porta SMTP inválida: %w", err)
	}

	service := &EmailService{
		Host:      os.Getenv("SMTP_HOST"),
		Port:      port,
		Username:  os.Getenv("SMTP_USER"),
		Password:  os.Getenv("SMTP_PASS"),
		FromName:  os.Getenv("SMTP_FROM_NAME"),
		FromEmail: os.Getenv("SMTP_FROM_EMAIL"),
	}

	service.auth = smtp.PlainAuth("", service.Username, service.Password, service.Host)
	return service, nil
}

// EmailData contém os dados para um template de e-mail.
type EmailData struct {
	ToName      string
	ToEmail     string
	Subject     string
	Body        string // Para e-mails de texto simples
	Template    string // Nome do arquivo de template HTML
	TemplateData interface{} // Dados a serem injetados no template
}

// Send envia um e-mail usando um template HTML ou texto simples.
func (s *EmailService) Send(data EmailData) error {
	var bodyContent bytes.Buffer
	var mimeHeaders string

	if data.Template != "" {
		// Usa template HTML
		t, err := template.ParseFiles(fmt.Sprintf("internal/email/templates/%s", data.Template))
		if err != nil {
			return fmt.Errorf("falha ao carregar o template de e-mail: %w", err)
		}
		t.Execute(&bodyContent, data.TemplateData)
		mimeHeaders = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	} else {
		// Usa texto simples
		bodyContent.WriteString(data.Body)
		mimeHeaders = "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	}

	headers := fmt.Sprintf("From: %s <%s>\r\nTo: %s <%s>\r\nSubject: %s\r\n%s",
		s.FromName, s.FromEmail, data.ToName, data.ToEmail, data.Subject, mimeHeaders)

	msg := append([]byte(headers), bodyContent.Bytes()...)
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	err := smtp.SendMail(addr, s.auth, s.FromEmail, []string{data.ToEmail}, msg)
	if err != nil {
		return fmt.Errorf("falha ao enviar e-mail: %w", err)
	}

	log.Printf("E-mail enviado com sucesso para %s com o assunto: %s", data.ToEmail, data.Subject)
	return nil
}
