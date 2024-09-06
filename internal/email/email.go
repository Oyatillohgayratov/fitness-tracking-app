package email

import "net/smtp"

func SendResetEmail(email, token string) error {
	from := "dilshoddilmurodov112@gmail.com"
	password := "xmxu rdhp pmdf pezk"
	to := email
	subject := "Password Reset Request"
	body := "your token: " + token

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	return err
}
