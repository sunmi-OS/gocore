package gomail

import (
	"github.com/go-gomail/gomail"
	"BITU-service/core/viper"
)

var mail *gomail.Dialer

func linkService() {

	mail = gomail.NewDialer(viper.C.GetString("email.host"), viper.C.GetInt("email.port"), viper.C.GetString("email.username"), viper.C.GetString("email.password"))
}

func SendRegisterEmail(email string, code string) error {
	if mail == nil {
		linkService()
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", "server@bitu.io", "BITU")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Welcome to register BITU users(欢迎注册BITU用户)")
	m.SetBody("text/html", "Your verification code is "+code+" (您的验证码是"+code+")  有效期为30分钟")
	return mail.DialAndSend(m)
}

func SendapplyResetPwdEmail(email string, code string) error {
	if mail == nil {
		linkService()
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", "server@bitu.io", "BITU")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Welcome to register BITU users(BITU用户找回密码)")
	m.SetBody("text/html", "Your verification code is "+code+" (您的验证码是"+code+")  有效期为30分钟")
	return mail.DialAndSend(m)
}
