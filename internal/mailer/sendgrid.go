package mailer

import (
	"bytes"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"log"
	"time"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {

	if fromEmail == "" {
		fmt.Println("Warning: fromEmail is empty. Ensure the FROM_EMAIL environment variable is set.")
	} else {
		fmt.Printf("fromEmail: %s\n", fromEmail)
	}

	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {

	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// template parse

	tmpl, err := template.ParseFS(FS, "template/"+templateFile)
	/*
		note:-

			The template.ParseFS function parses the plain text
			content into a *template.Template object.
			This object is used to render the template
			by replacing placeholders (e.g., {{.Username}}) with actual data.

			Yes, the embedded content is accessed from the binary.
			When you use the embed package, the specified files
			(e.g., templates) are embedded into the compiled binary
			during the build process. At runtime, the embed.FS provides
			a virtual file system interface to access this content.
	*/
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	/*
		note:-
			Template Execution: It executes the "subject" section of the parsed template
			(tmpl), which is defined in the user_invitation.tmpl file as:


			{{define "subject"}}Finish Registration with GoSocial{{end}}
			Data Binding: The data parameter is passed to the template.
			If the template contains placeholders (e.g., {{.Username}}), they are replaced with corresponding values from data.


			Output to Buffer: The rendered output (in this case, "Finish Registration with GoSocial")
			is written to the subject buffer, which is a bytes.Buffer.
	*/
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

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
