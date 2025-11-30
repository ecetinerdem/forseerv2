package mailer

const (
	fromName = "ForSeer"
)

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}
