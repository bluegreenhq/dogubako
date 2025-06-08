package adapter

import "context"

type MailAdapter interface {
	SendMail(ctx context.Context, from, to string, subject, body string) error
}
