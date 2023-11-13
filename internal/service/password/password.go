package password

import (
	"fmt"
	"password-self-service/internal/service/captcha"
	"password-self-service/pkg/gredis"
	"password-self-service/pkg/ldap"
	"password-self-service/pkg/logging"
)

var Service service

type service struct{}

// ResetPassword 重置密码
func (s *service) ResetPassword(username, newPassword, code string) error {
	user, err := captcha.Service.VerifyCaptcha(username, "password", code)
	if err != nil {
		return err
	}

	client := ldap.NewLdapClient()
	err = client.Connect()
	if err != nil {
		return err
	}
	defer client.Conn.Close()

	err = client.ModifyPassword(user.DN, newPassword)
	if err != nil {
		return err
	}

	gredis.Client.Del(gredis.Ctx, fmt.Sprintf("%s:%s", "password", user.Username))

	logging.Logger().Sugar().Infof("用户%v重置密码成功.", user.Username)

	return nil
}
