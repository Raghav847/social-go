package mailer

import "embed"

const (
	FromName            = ""
	maxRetries          = ""
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any) (int, error)
}
