package mailer

import (
	"bytes"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"time"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) send(templateFile, username, email string, data any, isSandbox bool) error {

	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	subject := new(bytes.Buffer)
	body := new(bytes.Buffer)
	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	for i := 0; i < maxRetries; i++ {
		response, err := m.client.Send(message)
		if err != nil {
			log.Printf("Failed to send email: %v, attempt %d of %d", err, i+1, maxRetries)
			log.Printf("Error : %v", err)

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		log.Printf("Successfully sent email: %v", response.StatusCode)
		return nil
	}

	return fmt.Errorf("failed to send email after %d attempts", maxRetries)
}
