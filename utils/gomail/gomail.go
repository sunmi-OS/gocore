package gomail

import (
	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"gopkg.in/gomail.v2"
)

var mail *gomail.Dialer

func linkService() {
	mail = gomail.NewDialer(
		viper.GetEnvConfig("email.host").String(),
		viper.GetEnvConfig("email.port").Int(),
		viper.GetEnvConfig("email.username").String(),
		viper.GetEnvConfig("email.password").String(),
	)
}

func SendEmail(email, fromMail, formNmae, subject, text string) error {
	if mail == nil {
		linkService()
	}
	m := gomail.NewMessage()
	m.SetAddressHeader("From", fromMail, formNmae)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)
	return mail.DialAndSend(m)
}
