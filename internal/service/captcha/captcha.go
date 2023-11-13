package captcha

import (
	"errors"
	"fmt"
	"math/rand"
	"password-self-service/pkg/config"
	"password-self-service/pkg/gredis"
	"password-self-service/pkg/ldap"
	"password-self-service/pkg/logging"
	"password-self-service/pkg/mail"
	"time"
)

var Service service

type service struct{}

// SendCaptcha 获取验证码
func (s *service) SendCaptcha(username, category string) error {
	client := ldap.NewLdapClient()
	err := client.Connect()
	if err != nil {
		return err
	}

	user, err := client.Search(username)
	if err != nil {
		return err
	}

	if user.Email == "" {
		return errors.New("该用户未配置邮箱,请联系管理员检查LDAP用户信息")
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	// 把验证码信息放到 redis，以便于验证时拿到
	gredis.Client.Set(gredis.Ctx, fmt.Sprintf("%s:%s", category, user.Username), code, 5*time.Minute)

	// 根据category拼接发送主题
	var title string
	if category == "password" {
		title = "密码重置"
	} else if category == "account" {
		title = "帐号解锁"
	}

	subject := fmt.Sprintf("%s-%s-验证码", config.Mail.From, title)
	content := fmt.Sprintf(`<div>
        <div>
            %s, 您好!
        </div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>您的AD域控帐号(%s)正在进行%s, 本次的验证码为 <span style="color: green;">%s</span> ,为了保证账号安全，验证码有效期为<span style="color: orange;">5分钟</span>。</p>
            <p>请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
        </div>
        <div>
            <p>此邮箱为系统邮箱，请勿回复。</p>
        </div>
    </div>`, user.Nickname, user.Username, title, code)

	// 发送验证码消息
	err = mail.SendMail([]string{user.Email}, subject, content)
	if err != nil {
		return err
	}

	logging.Logger().Sugar().Infof("用户%v进行%v发送了验证码.", user.Username, title)

	return nil
}

// VerifyCaptcha 校验验证码
func (s *service) VerifyCaptcha(username, category, code string) (*ldap.Attributes, error) {
	values, _ := gredis.Client.Get(gredis.Ctx, fmt.Sprintf("%s:%s", category, username)).Result()

	if values == "" {
		return nil, errors.New("验证码已过期，请重新获取")
	}

	if values != code {
		return nil, errors.New("验证码错误，请重新输入")
	}

	client := ldap.NewLdapClient()
	err := client.Connect()
	if err != nil {
		return nil, err
	}
	defer client.Conn.Close()

	user, err := client.Search(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
