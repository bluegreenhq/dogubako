package adapter

import (
	"context"

	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

type SMTPMailAdapter struct {
	host     string
	port     int
	user     string
	password string
}

var _ MailAdapter = (*SMTPMailAdapter)(nil)

func NewSMTPMailAdapter(host string, port int, user, password string) (MailAdapter, error) {
	return &SMTPMailAdapter{
		host:     host,
		port:     port,
		user:     user,
		password: password,
	}, nil
}

func (a SMTPMailAdapter) SendMail(ctx context.Context, from, to string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(a.host, a.port, a.user, a.password)
	if err := d.DialAndSend(m); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
