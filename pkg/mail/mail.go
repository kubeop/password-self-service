package mail

import (
	"crypto/tls"
	"password-self-service/pkg/config"

	"gopkg.in/gomail.v2"
)

func SendMail(mailTo []string, subject string, content string) error {
	newMail := gomail.NewMessage()

	newMail.SetHeader("From", newMail.FormatAddress(config.Mail.User, config.Mail.From))

	// 发送给多个用户
	newMail.SetHeader("To", mailTo...)
	// 设置邮件主题
	newMail.SetHeader("Subject", subject)
	// 设置邮件正文
	newMail.SetBody("text/html", content)

	do := gomail.NewDialer(config.Mail.Host, config.Mail.Port, config.Mail.User, config.Mail.Password)
	if config.Mail.TLS {
		do.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return do.DialAndSend(newMail)
}
