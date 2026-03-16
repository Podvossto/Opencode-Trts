// Purpose: Email service — transactional email dispatch (SMTP via gomail)
// Ref: API Contract — POST /api/v1/applications/{id}/email
// Dev must implement: Send(to, subject, body string) error, template rendering
package email

import (
	"fmt"
	"log"

	gomail "gopkg.in/gomail.v2"
)

// Config holds SMTP configuration
type Config struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

// Sender defines the email sending contract
type Sender interface {
	Send(to, subject, htmlBody string) error
}

// SMTPSender implements Sender using SMTP
type SMTPSender struct {
	cfg    Config
	dialer *gomail.Dialer
}

// NewSMTPSender creates a new SMTPSender
func NewSMTPSender(cfg Config) *SMTPSender {
	return &SMTPSender{
		cfg:    cfg,
		dialer: gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Pass),
	}
}

// Send dispatches an email to the given recipient
func (s *SMTPSender) Send(to, subject, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("email send failed: %w", err)
	}
	log.Printf("[email] sent to %s: %s", to, subject)
	return nil
}
