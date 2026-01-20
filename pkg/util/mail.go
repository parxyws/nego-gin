package util

import (
	"bytes"
	"path"
	"text/template"

	"github.com/parxyws/nego-gin/config"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func GenerateOTPMailMessage(cfg *config.Config, user *domain.User, otp string) (*gomail.Message, error) {
	type OTPMessage struct {
		Name string
		OTP  string
	}

	otpData := &OTPMessage{
		Name: user.FirstName + " " + user.LastName,
		OTP:  otp,
	}

	logrus.WithFields(logrus.Fields{"function": "GenerateOTPMailMessage"}).Debugf("first=%q last=%q", user.FirstName, user.LastName)
	filePath := path.Join("web", "template", "otp.html")

	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, otpData); err != nil {
		return nil, err
	}

	emailMsg := gomail.NewMessage()
	emailMsg.SetHeader("From", cfg.Mail.User)
	emailMsg.SetHeader("To", user.Email)
	emailMsg.SetHeader("Subject", "Your OTP Verification Code")
	emailMsg.SetBody("text/html", body.String())

	return emailMsg, nil
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

	emailMsg := gomail.NewMessage()
	emailMsg.SetHeader("From", cfg.Mail.User)
	emailMsg.SetHeader("To", to)
	emailMsg.SetHeader("Subject", "Password Rest Requested")
	emailMsg.SetBody("text/html", body.String())

	return emailMsg, nil
}
