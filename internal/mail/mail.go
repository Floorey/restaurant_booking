package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

type Mailer struct {
	auth smtp.Auth
	host string
	port string
	from string
}

func New() *Mailer {
	return &Mailer{
		auth: smtp.PlainAuth("", os.Getenv("SMTP_USER"),
			os.Getenv("SMTP_PASS"), os.Getenv("SMTP_HOST")),
		host: os.Getenv("SMTP_HOST"),
		port: os.Getenv("SMTP_PORT"),
		from: os.Getenv("SMTP_ADDR"),
	}
}

func (m *Mailer) Send(to, subj, html string) error {
	e := email.NewEmail()
	e.From = m.from
	e.To = []string{to}
	e.Subject = subj
	e.HTML = []byte(html)
	return e.send(fmt.Sprintf("%s:%s", m.host, m.port), m.auth)
}
