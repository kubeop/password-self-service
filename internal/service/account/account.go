package account

import (
	"fmt"
	"password-self-service/internal/service/captcha"
	"password-self-service/pkg/gredis"
	"password-self-service/pkg/ldap"
	"password-self-service/pkg/logging"
)

var Service service

type service struct{}

// UnlockAccount 解锁账号
func (s *service) UnlockAccount(username, code string) error {
	user, err := captcha.Service.VerifyCaptcha(username, "account", code)
	if err != nil {
		return err
	}

	client := ldap.NewLdapClient()
	err = client.Connect()
	if err != nil {
		return err
	}
	defer client.Conn.Close()

	err = client.UnlockAccount(user.DN)
	if err != nil {
		return err
	}

	gredis.Client.Del(gredis.Ctx, fmt.Sprintf("%s:%s", "account", user.Username))

	logging.Logger().Sugar().Infof("用户%v解锁账号成功.", user.Username)

	return nil
}
