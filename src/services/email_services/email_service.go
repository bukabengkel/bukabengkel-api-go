package email_services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/peang/bukabengkel-api-go/src/config"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type EmailAction string

const (
	WAITING_FOR_PAYMENT EmailAction = "order_distributor/waiting_for_payment.html"
	ORDER_PAID          EmailAction = "order_distributor/order_paid.html"
)

type WaitingForPaymentData struct {
	StoreName string
	ExpiredAt string
	PaymentURL string
}

type EmailService interface {
	SendWaitingForPaymentEmail(ctx context.Context, to string, data WaitingForPaymentData) error
}

const (
	EMAIL_SERVICE_MAILTRAP = "mailtrap"
	EMAIL_SERVICE_RESEND   = "resend"
)

func NewEmailService(config *config.Config, logger utils.Logger) (EmailService, error) {
	if config.Env != "production" {
		return newMailtrapService(config, logger), nil
	}

	switch config.EmailProvider.EmailServiceName {
	case EMAIL_SERVICE_RESEND:
		return newResendService(config), nil
	default:
		return nil, fmt.Errorf("email service %s not supported", config.EmailProvider.EmailServiceName)
	}
}

func RenderTemplate(templateName string, data map[string]any) (string, error) {
	// Build template file path - templateName should include the .html extension
	templatePath := filepath.Join("src/services/email_services/templates", templateName)

	// Check if template file exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template file not found: %s", templatePath)
	}

	// Read template file
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %v", err)
	}

	// Parse and execute template using Go template syntax
	tmpl, err := template.New(templateName).Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return result.String(), nil
}

// Helper function to get subject based on action
func getSubjectForAction(action string) string {
	switch action {
	case "order_distributor/waiting_for_payment.html":
		return "Menunggu Pembayaran - BukaBengkel"
	default:
		return "BukaBengkel Notification"
	}
}

// Helper function to format deadline time
func FormatDeadlineForEmail(deadline time.Time) string {
	months := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	
	return fmt.Sprintf("%02d:%02d %d %s %d", 
		deadline.Hour(), 
		deadline.Minute(), 
		deadline.Day(), 
		months[deadline.Month()], 
		deadline.Year())
}
