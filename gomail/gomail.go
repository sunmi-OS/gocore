package gomail

import (
	"github.com/go-gomail/gomail"
	"github.com/sunmi-OS/gocore/viper"
)

var mail *gomail.Dialer

func linkService() {

	 mail = gomail.NewDialer(viper.C.GetString("email.host"), viper.C.GetInt("email.port"), viper.C.GetString("email.username"), viper.C.GetString("email.password"))
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
