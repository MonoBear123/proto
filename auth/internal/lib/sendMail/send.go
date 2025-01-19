package sendMail

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

var (
	VariotyBody = map[string]string{
		"Активация аккаунта": `<!DOCTYPE html>
  <html>
  <head>
    <title>Email Template</title>
  </head>
  <body>
    <p>Здравствуйте!</p>
    <p>Спасибо за регистрацию на нашем сервисе. Для активации вашего аккаунта, пожалуйста, перейдите по следующей <a href="http://%s/activate?token=%s">ссылке</a>.</p>
    <p>Команда вашего сервиса.</p>
  </body>
  </html>`,
		"Сброс пароля": `<!DOCTYPE html>
  <html>
  <head>
    <title>Email Template</title>
  </head>
  <body>
    <p>Здравствуйте!</p>
    <p>Если вы забыли пароль, вы можете сбросить его, перейдя по следующей <a href="http://%s/reset-password?token=%s">ссылке</a>.</p>
    <p>Команда вашего сервиса.</p>
  </body>
  </html>`,
	}
)

func SendMessagee(email, subject, token string) error {
	body := CreateBody(subject, token)
	msg := gomail.NewMessage()
	msg.SetHeader("From", "predict.service@mail.ru")
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	n := gomail.NewDialer("smtp.mail.ru", 587, "predict.service@mail.ru", "2KRFNX49efzQ9z5r4cjQ")

	if err := n.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}

func CreateBody(subject, token string) string {
	link := "185.112.102.43:8080"
	body := VariotyBody[subject]
	return fmt.Sprintf(body, link, token)

}
