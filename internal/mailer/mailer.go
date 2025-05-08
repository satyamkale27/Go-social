package mailer

const (
	FromName   = "Gosocial"
	maxRetries = 3
)

type Client interface {
	send(templateFile, username, email string, data any, isSandbox bool) error
}
