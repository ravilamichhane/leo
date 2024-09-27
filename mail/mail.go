package mail

import (
	"crypto/tls"
	"embed"
	"fmt"
	"net/smtp"

	"github.com/ravilmc/leo/web"

	"github.com/ravilmc/leo/email"

	"github.com/google/uuid"
)

type Mail struct {
	To           []string
	Subject      string
	Bcc          []string
	Cc           []string
	Template     string
	Data         map[string]interface{}
	Body         []byte
	XEntityRefID *uuid.UUID
}

//go:embed templates/*
var templates embed.FS

func SendMail(mail Mail) error {

	go func() {
		host := web.GetEnv("MAIL_HOST", "")
		port := web.GetEnv("MAIL_PORT", "")
		user := web.GetEnv("MAIL_USER", "")
		password := web.GetEnv("MAIL_PASSWORD", "")
		tlsEnabled := web.GetEnvBool("MAIL_TLS", false)
		address := host + ":" + port

		e := email.NewEmail()
		e.From = fmt.Sprintf("Big Bracket Esports <%s>", user)
		e.To = mail.To
		e.Bcc = mail.Bcc
		e.Cc = mail.Cc
		e.Subject = mail.Subject

		if mail.XEntityRefID != nil {
			e.Headers.Add("X-Entity-Ref-ID", mail.XEntityRefID.String())
		}

		// body, err := template.ParseFS(templates, "templates/"+mail.Template)

		// var bodyBuffer bytes.Buffer
		// err = body.Execute(&bodyBuffer, mail.Data)

		e.HTML = mail.Body

		auth := smtp.PlainAuth("", user, password, host)

		if tlsEnabled {
			e.SendWithTLS(address, auth, &tls.Config{
				InsecureSkipVerify: false,
				ServerName:         host,
			})
		} else {
			e.Send(address, auth)
		}

	}()

	return nil
}
