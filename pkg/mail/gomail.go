package mail

import (
	"github.com/parxyws/nego-gin/config"
	"gopkg.in/gomail.v2"
)

func NewGoMailDialer(cfg *config.Config) *gomail.Dialer {
	return gomail.NewDialer(cfg.Mail.Host, cfg.Mail.Port, cfg.Mail.User, cfg.Mail.Password)
}
