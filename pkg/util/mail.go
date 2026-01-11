package util

import (
	"bytes"
	"path"
	"text/template"

	"github.com/parxyws/nego-gin/config"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"gopkg.in/gomail.v2"
)

func GenerateOTPMailMessage(cfg *config.Config, user *domain.User, otp string) (*gomail.Message, error) {
	type OTPMessage struct {
		Name string
		OTP  string
	}

	otpData := &OTPMessage{
		Name: user.FirstName,
		OTP:  otp,
	}
	filePath := path.Join("web", "template", "otp.html")

	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, otpData); err != nil {
		return nil, err
	}

	mailMessage := gomail.NewMessage()
	mailMessage.SetHeader("From", cfg.Mail.User)
	mailMessage.SetHeader("To", user.Email)
	mailMessage.SetHeader("Subject", "Your OTP Verification Code")
	mailMessage.SetBody("text/html", body.String())

	return mailMessage, nil
}

func GenerateLinkMailMessage(cfg *config.Config, to string, link string) (*gomail.Message, error) {
	type PasswordResetMessage struct {
		Link string
	}

	resetData := &PasswordResetMessage{
		Link: link,
	}
	filePath := path.Join("web", "template", "reset-password.html")
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, resetData); err != nil {
		return nil, err
	}

	mailMessage := gomail.NewMessage()
	mailMessage.SetHeader("From", cfg.Mail.User)
	mailMessage.SetHeader("To", to)
	mailMessage.SetHeader("Subject", "Password Rest Requested")
	mailMessage.SetBody("text/html", body.String())

	return mailMessage, nil
}
