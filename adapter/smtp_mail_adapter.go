package adapter

import (
	"github.com/pkg/errors"
	"github.com/wneessen/go-mail"
)

type SMTPMailAdapter struct {
	client *mail.Client
}

var _ MailAdapter = (*SMTPMailAdapter)(nil)

func NewSMTPMailAdapter(host string, port int, user, password string) (MailAdapter, error) {
	client, err := mail.NewClient(
		host,
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(user),
		mail.WithPassword(password),
		mail.WithPort(port),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &SMTPMailAdapter{
		client: client,
	}, nil
}

func (a SMTPMailAdapter) SendMail(from, to string, subject, body string) error {
	m := mail.NewMsg()

	err := m.From(from)
	if err != nil {
		return errors.WithStack(err)
	}

	err = m.To(to)
	if err != nil {
		return errors.WithStack(err)
	}

	m.Subject(subject)
	m.SetBodyString(mail.TypeTextPlain, body)

	if err := a.client.DialAndSend(m); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
