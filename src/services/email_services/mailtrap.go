package email_services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/peang/bukabengkel-api-go/src/config"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type MailtrapService struct {
	logger utils.Logger
	Key    string
}

type MailtrapEmailRequest struct {
	From     MailtrapSender      `json:"from"`
	To       []MailtrapRecipient `json:"to"`
	Subject  string              `json:"subject"`
	HTML     string              `json:"html,omitempty"`
	Text     string              `json:"text,omitempty"`
	Category string              `json:"category,omitempty"`
}

type MailtrapSender struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type MailtrapRecipient struct {
	Email string `json:"email"`
}

func newMailtrapService(config *config.Config, logger utils.Logger) *MailtrapService {
	return &MailtrapService{Key: config.EmailProvider.EmailAPIKey, logger: logger}
}

func (s *MailtrapService) SendWaitingForPaymentEmail(ctx context.Context, to string, data WaitingForPaymentData) error {
	// Render the HTML template
	htmlContent, err := RenderTemplate("order_distributor/waiting_for_payment.html", map[string]any{
		"StoreName":  data.StoreName,
		"ExpiredAt":  data.ExpiredAt,
		"PaymentURL": data.PaymentURL,
	})
	if err != nil {
		return fmt.Errorf("failed to render template: %v", err)
	}

	return s.sendEmail(ctx, to, "Menunggu Pembayaran - BukaBengkel", htmlContent)
}

func (s *MailtrapService) sendEmail(ctx context.Context, to string, subject string, htmlContent string) error {
	url := "https://sandbox.api.mailtrap.io/api/send/2406850"

	// Create the email request payload
	emailRequest := MailtrapEmailRequest{
		From: MailtrapSender{
			Email: "no-reply@bukabengkel.id",
			Name:  "Buka Bengkel",
		},
		To: []MailtrapRecipient{
			{
				Email: to,
			},
		},
		Subject: subject,
		HTML:    htmlContent,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(emailRequest)
	if err != nil {
		s.logger.Error("failed to marshal email request", "error", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		s.logger.Error("failed to create request", "error", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.Key)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("failed to send email", "error", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("failed to read response body", "error", err)
	}
	fmt.Println(string(body))
	fmt.Println(resp.StatusCode)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("mailtrap API returned status", "status", resp.StatusCode)
	}

	s.logger.Info("email sent successfully", "status", resp.StatusCode)
	return nil
}
