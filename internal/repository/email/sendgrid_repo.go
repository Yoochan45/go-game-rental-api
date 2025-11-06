package email

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

var (
	htmlTagRegex = regexp.MustCompile(`<[^>]*>`)
	emailRegex   = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type EmailRepository interface {
	SendEmail(ctx context.Context, to, subject, plainText, htmlContent string) error
	SendWithTemplate(ctx context.Context, to, templateID string, dynamicData map[string]interface{}) error
}

type SendGridRepository struct {
	client   *sendgrid.Client
	fromName string
	fromAddr string
}

func NewSendGridRepository() (*SendGridRepository, error) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	fromAddr := os.Getenv("SENDGRID_FROM_EMAIL")
	fromName := os.Getenv("SENDGRID_FROM_NAME")

	if apiKey == "" || fromAddr == "" {
		return nil, fmt.Errorf("sendgrid not configured: missing API_KEY or FROM_EMAIL")
	}
	if fromName == "" {
		fromName = "Game Rental"
	}

	return &SendGridRepository{
		client:   sendgrid.NewSendClient(apiKey),
		fromName: fromName,
		fromAddr: fromAddr,
	}, nil
}

func (s *SendGridRepository) SendEmail(ctx context.Context, to, subject, plainText, htmlContent string) error {
	_ = ctx // ctx unused - SendGrid client doesn't support context timeout
	if !isValidEmail(to) {
		return fmt.Errorf("invalid email address: %s", to)
	}

	// Fallback plainText from HTML if empty to avoid spam marking
	if plainText == "" && htmlContent != "" {
		plainText = stripHTML(htmlContent)
	}

	from := mail.NewEmail(s.fromName, s.fromAddr)
	toEmail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(from, subject, toEmail, plainText, htmlContent)

	resp, err := s.client.Send(message)
	if err != nil {
		logrus.WithError(err).WithField("to", to).Error("SendGrid send failed")
		return fmt.Errorf("failed to send email: %w", err)
	}
	if resp.StatusCode >= 400 {
		logrus.WithFields(logrus.Fields{
			"status": resp.StatusCode,
			"body":   resp.Body,
			"to":     to,
		}).Error("SendGrid error")
		return fmt.Errorf("sendgrid error: status=%d", resp.StatusCode)
	}
	logrus.WithFields(logrus.Fields{
		"to":      to,
		"subject": subject,
	}).Info("Email sent successfully")
	return nil
}

func (s *SendGridRepository) SendWithTemplate(ctx context.Context, to, templateID string, dynamicData map[string]interface{}) error {
	_ = ctx // ctx unused - SendGrid client doesn't support context timeout
	if s.client == nil || s.fromAddr == "" {
		return fmt.Errorf("sendgrid not configured")
	}
	from := mail.NewEmail(s.fromName, s.fromAddr)
	toEmail := mail.NewEmail("", to)

	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.SetTemplateID(templateID)

	p := mail.NewPersonalization()
	p.AddTos(toEmail)
	for k, v := range dynamicData {
		p.SetDynamicTemplateData(k, v)
	}
	message.AddPersonalizations(p)

	resp, err := s.client.Send(message)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"to":         to,
			"template_id": templateID,
		}).Error("SendGrid template send failed")
		return fmt.Errorf("failed to send template email: %w", err)
	}
	if resp.StatusCode >= 400 {
		logrus.WithFields(logrus.Fields{
			"status":     resp.StatusCode,
			"to":         to,
			"template_id": templateID,
		}).Error("SendGrid template error")
		return fmt.Errorf("sendgrid template error: status=%d", resp.StatusCode)
	}
	logrus.WithFields(logrus.Fields{
		"to":         to,
		"template_id": templateID,
	}).Info("Template email sent successfully")
	return nil
}

type MockEmailRepository struct {
	SentEmails []MockEmail
}

type MockEmail struct {
	To          string
	Subject     string
	PlainText   string
	HTMLContent string
	TemplateID  string
	Data        map[string]interface{}
}

func (m *MockEmailRepository) SendEmail(ctx context.Context, to, subject, plainText, htmlContent string) error {
	m.SentEmails = append(m.SentEmails, MockEmail{
		To:          to,
		Subject:     subject,
		PlainText:   plainText,
		HTMLContent: htmlContent,
	})
	return nil
}

func (m *MockEmailRepository) SendWithTemplate(ctx context.Context, to, templateID string, dynamicData map[string]interface{}) error {
	m.SentEmails = append(m.SentEmails, MockEmail{
		To:         to,
		TemplateID: templateID,
		Data:       dynamicData,
	})
	return nil
}

// stripHTML removes HTML tags for plaintext fallback
func stripHTML(html string) string {
	plain := htmlTagRegex.ReplaceAllString(html, "")
	return strings.TrimSpace(plain)
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}