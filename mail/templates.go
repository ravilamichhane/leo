package mail

import (
	"fmt"

	"github.com/ravilmc/leo/web"

	"github.com/google/uuid"
	"github.com/matcornic/hermes/v2"
)

type RegisterMailContent struct {
	Name string
	Link string
	OTP  string
}

var h = hermes.Hermes{
	Theme: &hermes.Flat{},
	Product: hermes.Product{
		Name:      "Big Bracket Esports",
		Link:      "http://esports.bigbracket.io",
		Logo:      "https://www.bigbracket.io/logo.png",
		Copyright: "Copyright © 2024 Big Bracket Esports. All rights reserved.",
	},
}

func SendRegisterMail(to string, data RegisterMailContent) error {
	URL := web.GetEnv("GOOGLE_REDIRECT_URL", "http://localhost:5173")
	email := hermes.Email{

		Body: hermes.Body{
			Name: data.Name,
			Intros: []string{
				"Welcome to Big Bracket Esports! We're very excited to have you on board.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Use the code below to verify your account",
					InviteCode:   data.OTP,
				},
				{
					Instructions: "Or click on the link below",
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Confirm your account",
						Link:  fmt.Sprintf("%s/verify-email?token=%s&email=%s", URL, data.OTP, to),
					},
				},
			},
		},
	}

	body, err := h.GenerateHTML(email)
	if err != nil {
		return err
	}
	uid := uuid.New()

	mail := Mail{
		To:           []string{to},
		Subject:      "Register",
		Template:     "new_register.tmpl",
		Body:         []byte(body),
		XEntityRefID: &uid,
	}
	return SendMail(mail)
}

type PasswordResetMailContent struct {
	Name string
	Link string
	OTP  string
}

func SendPasswordResetMail(to string, data PasswordResetMailContent) error {
	mail := Mail{
		To:       []string{to},
		Subject:  "Reset Password",
		Template: "forgot_password.tmpl",
		Data: map[string]interface{}{
			"Name": data.Name,
			"Link": data.Link,
			"OTP":  data.OTP,
		},
	}
	return SendMail(mail)
}
