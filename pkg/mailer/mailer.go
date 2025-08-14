package mailer

import (
	"fmt"
	"net/smtp"
	"strings"
)

// Mailer interface for sending emails
type Mailer interface {
	SendEmail(to, subject, body string) error
	SendVerificationEmail(to, token string) error
	SendPasswordResetEmail(to, token string) error
}

// SMTPMailer implements Mailer interface using SMTP
type SMTPMailer struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// NewSMTPMailer creates a new SMTP mailer
func NewSMTPMailer(host string, port int, username, password, from string) *SMTPMailer {
	return &SMTPMailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// SendEmail sends a plain text email
func (m *SMTPMailer) SendEmail(to, subject, body string) error {
	// Create message
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", m.from, to, subject, body)

	// Connect to SMTP server
	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	// Send email
	err := smtp.SendMail(addr, auth, m.from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendVerificationEmail sends an email verification email
func (m *SMTPMailer) SendVerificationEmail(to, token string) error {
	subject := "Verify your email address"
	body := fmt.Sprintf(`
Hello!

Please verify your email address by clicking the following link:

https://4planet.local/verify-email?token=%s

This link will expire in 24 hours.

If you didn't create an account, please ignore this email.

Best regards,
4Planet Team
`, token)

	return m.SendEmail(to, subject, strings.TrimSpace(body))
}

// SendPasswordResetEmail sends a password reset email
func (m *SMTPMailer) SendPasswordResetEmail(to, token string) error {
	subject := "Reset your password"
	body := fmt.Sprintf(`
Hello!

You requested to reset your password. Click the following link to set a new password:

https://4planet.local/reset-password?token=%s

This link will expire in 1 hour.

If you didn't request a password reset, please ignore this email.

Best regards,
4Planet Team
`, token)

	return m.SendEmail(to, subject, strings.TrimSpace(body))
}

// NoOpMailer is a mock mailer for development/testing
type NoOpMailer struct{}

// NewNoOpMailer creates a new no-op mailer
func NewNoOpMailer() *NoOpMailer {
	return &NoOpMailer{}
}

// SendEmail does nothing (for development)
func (m *NoOpMailer) SendEmail(to, subject, body string) error {
	// Log the email for development purposes
	fmt.Printf("[MAILER] Would send email to %s\nSubject: %s\nBody: %s\n", to, subject, body)
	return nil
}

// SendVerificationEmail does nothing (for development)
func (m *NoOpMailer) SendVerificationEmail(to, token string) error {
	fmt.Printf("[MAILER] Would send verification email to %s with token %s\n", to, token)
	return nil
}

// SendPasswordResetEmail does nothing (for development)
func (m *NoOpMailer) SendPasswordResetEmail(to, token string) error {
	fmt.Printf("[MAILER] Would send password reset email to %s with token %s\n", to, token)
	return nil
}
