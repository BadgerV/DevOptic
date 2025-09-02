// package emailservice

// import (
// 	"fmt"
// 	"log"
// 	"net/smtp"
// 	"os"
// 	"strings"
// )

// // Email Service struct hold SMTP config
// type EmailService struct {
// 	Host     string
// 	Port     string
// 	Username string
// 	Password string
// 	From     string
// 	To       []string
// }

// func NewEmailService() *EmailService {
// 	return &EmailService{
// 		Host:     os.Getenv("SMTP_HOST"),
// 		Port:     os.Getenv("SMTP_PORT"),
// 		Username: os.Getenv("SMTP_USERNAME"),
// 		Password: os.Getenv("SMTP_PASSWORD"),
// 		From:     os.Getenv("SMTP_EMAIL_FROM"),
// 	}
// }

// // Sends an email with the provided subject, body and recipients
// func (e *EmailService) Send(subject, body string, to []string) error {
// 	println(e.Host)
// 	println(e.Username)
// 	println(e.Password)
// 	println(e.Port)
// 	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
// 	// auth := smtp.CRAMMD5Auth(e.Username, e.Password)

// 	msg := []byte(fmt.Sprintf(
// 		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
// 		e.From,
// 		strings.Join(e.To, ","),
// 		subject,
// 		body,
// 	))

// 	address := fmt.Sprintf("%s:%s", e.Host, e.Port)
// 	err := smtp.SendMail(address, auth, e.From, to, msg)

// 	if err != nil {
// 		log.Fatal("Failed to send email: %w", err)
// 	} else {
// 		log.Printf("Email sent successfully for the subject: %w", subject)
// 	}

// 	return nil
// }

package emailservice

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/mail.v2"
)

// EmailService struct holds SMTP config
type EmailService struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewEmailService() *EmailService {
	port := 587 // default SMTP port
	if portStr := os.Getenv("SMTP_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	return &EmailService{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_EMAIL_FROM"),
	}
}

// Send sends an email with the provided subject, body and recipients
func (e *EmailService) Send(subject, body string, to []string) error {
	// Create message
	m := mail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Create dialer with authentication
	d := mail.NewDialer(e.Host, e.Port, e.Username, e.Password)
	
	// For Office365/Outlook, you might need to enable StartTLS
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// Send email
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully with subject: %s", subject)
	return nil
}

// SendHTML sends an HTML email with the provided subject, body and recipients
func (e *EmailService) SendHTML(subject, htmlBody string, to []string) error {
	m := mail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := mail.NewDialer(e.Host, e.Port, e.Username, e.Password)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send HTML email: %v", err)
		return err
	}

	log.Printf("HTML email sent successfully with subject: %s", subject)
	return nil
}

// SendWithAttachment sends an email with attachment
func (e *EmailService) SendWithAttachment(subject, body string, to []string, attachmentPath string) error {
	m := mail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	
	// Attach file
	m.Attach(attachmentPath)

	d := mail.NewDialer(e.Host, e.Port, e.Username, e.Password)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email with attachment: %v", err)
		return err
	}

	log.Printf("Email with attachment sent successfully with subject: %s", subject)
	return nil
}