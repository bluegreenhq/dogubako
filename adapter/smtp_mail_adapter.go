package adapter

import (
	"github.com/pkg/errors"
	"github.com/wneessen/go-mail"
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

	d, err := mail.NewClient(
		a.host,
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(a.user),
		mail.WithPassword(a.password),
		mail.WithPort(a.port),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := d.DialAndSend(m); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
