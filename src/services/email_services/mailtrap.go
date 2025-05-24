package email_services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/peang/bukabengkel-api-go/src/config"
)

type MailtrapService struct {
	// client *mailtrap.Client
}

func newMailtrapService(config *config.Config) *MailtrapService {
	return &MailtrapService{}
}

func (s *MailtrapService) SendEmail(ctx context.Context, to string, subject string, body string) error {
	url := "https://sandbox.api.mailtrap.io/api/send/2406850"
	method := "POST"

	payload := strings.NewReader(`{\"from\":{\"email\":\"hello@example.com\",\"name\":\"Mailtrap Test\"},\"to\":[{\"email\":\"bukabengkel.id@gmail.com\"}],\"subject\":\"You are awesome!\",\"text\":\"Congrats for sending test email with Mailtrap!\",\"category\":\"Integration Test\"}`)

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Authorization", "Bearer ****8a8e")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	bodyResponse, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(bodyResponse))

	return nil
}