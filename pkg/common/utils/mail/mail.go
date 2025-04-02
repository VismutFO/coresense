package mail

import (
	"gopkg.in/gomail.v2"
)

func sendEmailGomail(to, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@yourdomain.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/plain", "Click the link: http://localhost:8080/verify?token="+token)

	d := gomail.NewDialer("smtp.yourmail.com", 587, "your_email", "your_password")
	return d.DialAndSend(m)
}
