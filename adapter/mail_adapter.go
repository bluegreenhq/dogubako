package adapter

type MailAdapter interface {
	SendMail(from, to string, subject, body string) error
}
