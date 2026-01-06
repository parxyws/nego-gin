package util

import (
	"bytes"
	"path"
	"text/template"

	"github.com/parxyws/nego-gin/config"
	"gopkg.in/gomail.v2"
)

func GenerateOTPMailMessage(cfg *config.Config, to string, otp string) *gomail.Message {
	filePath := path.Join("web", "template", "otp.html")

	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return nil
	}

	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, "otp", otp); err != nil {
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.Mail.User)
	m.SetHeader("To", to)
	m.SetBody("text/html", body.String())

	return m
}
