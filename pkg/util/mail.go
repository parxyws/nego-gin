package util

import (
	"bytes"
	"path"
	"text/template"

	"github.com/parxyws/nego-gin/config"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"gopkg.in/gomail.v2"
)

func GenerateOTPMailMessage(cfg *config.Config, user *domain.User, otp string) *gomail.Message {
	type OTPMessage struct {
		username string
		otp      string
	}

	otpData := &OTPMessage{
		username: user.Username,
		otp:      otp,
	}
	filePath := path.Join("web", "template", "otp.html")

	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return nil
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, otpData); err != nil {
		return nil
	}

	mailMessage := gomail.NewMessage()
	mailMessage.SetHeader("From", cfg.Mail.User)
	mailMessage.SetHeader("To", user.Email)
	mailMessage.SetBody("text/html", body.String())

	return mailMessage
}

func GenerateLinkMailMessage(cfg *config.Config, to string, link string) *gomail.Message {
	filePath := path.Join("web", "template", "reset-password.html")
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return nil
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, link); err != nil {
		return nil
	}

	mailMessage := gomail.NewMessage()
	mailMessage.SetHeader("From", cfg.Mail.User)
	mailMessage.SetHeader("To", to)
	mailMessage.SetBody("text/html", body.String())

	return mailMessage
}
