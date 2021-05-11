package gomail

import (
	viper2 "github.com/sunmi-OS/gocore/conf/viper"
	"gopkg.in/gomail.v2"
)

var mail *gomail.Dialer

func linkService() {
	mail = gomail.NewDialer(viper2.C.GetString("email.host"), viper2.C.GetInt("email.port"), viper2.C.GetString("email.username"), viper2.C.GetString("email.password"))
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
