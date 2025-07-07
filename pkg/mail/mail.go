package mail

import (
	"crypto/tls"
	"password-self-service/pkg/config"

	"gopkg.in/gomail.v2"
)

func SendMail(mailTo []string, subject string, content string) error {
	newMail := gomail.NewMessage()

	newMail.SetHeader("From", newMail.FormatAddress(config.Setting.Channel.Mail.User, config.Setting.Channel.Mail.From))

	// 发送给多个用户
	newMail.SetHeader("To", mailTo...)
	// 设置邮件主题
	newMail.SetHeader("Subject", subject)
	// 设置邮件正文
	newMail.SetBody("text/html", content)

	do := gomail.NewDialer(config.Setting.Channel.Mail.Host, config.Setting.Channel.Mail.Port, config.Setting.Channel.Mail.User, config.Setting.Channel.Mail.Password)
	if config.Setting.Channel.Mail.TLS {
		do.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return do.DialAndSend(newMail)
}
