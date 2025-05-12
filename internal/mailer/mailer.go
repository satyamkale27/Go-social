package mailer

import "embed"

const (
	FromName            = "Gosocial"
	maxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

/*
the embed package embeds the template file's content (as plain text)
into the Go binary at compile time. This allows the application
to access the template as if it were a file, but it is stored within the compiled binary.
The template remains in its original text format and is parsed and executed at runtime.
*/

//go:embed "template"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
}
