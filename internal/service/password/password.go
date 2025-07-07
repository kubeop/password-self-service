package password

import (
	"errors"
	"fmt"
	"password-self-service/internal/service/captcha"
	"password-self-service/pkg/aliyun"
	"password-self-service/pkg/config"
	"password-self-service/pkg/gredis"
	"password-self-service/pkg/ldap"
	"password-self-service/pkg/logging"
	"password-self-service/pkg/mail"
	"password-self-service/pkg/tencent"
	"password-self-service/pkg/utils"
	"strconv"
	"time"
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

	logging.Logger().Sugar().Infof("用户 %v 重置密码成功.", user.Username)

	return nil
}

// PasswordExpired 密码过期通知
func (s *service) PasswordExpired() error {
	client := ldap.NewLdapClient()
	err := client.Connect()
	if err != nil {
		return err
	}

	defer client.Conn.Close()

	users, err := client.CheckPasswordExpired()
	if err != nil {
		logging.Logger().Sugar().Errorf("密码过期检查失败: %v", err)
		return err
	}

	for _, user := range users {
		PasswordExpireTime, err := strconv.ParseInt(user.PasswordExpire, 10, 64)
		if err != nil {
			logging.Logger().Error(err.Error())
			return err
		}

		if user.Mobile != "" && user.Email != "" {
			// 计算密码过期时间
			expirationTime := time.Unix(PasswordExpireTime/1e7-11644473600, 0)
			expireTime := expirationTime.AddDate(0, 0, 90)

			// 判定密码过期时间是否小于7天，小于7天则发送提醒修改密码消息，已过期不再提醒
			if time.Until(expireTime) < 7*24*time.Hour && time.Until(expireTime) >= 0 {
				// 发送修改密码消息
				logging.Logger().Sugar().Infof("用户：%v , 中文名：%v , 手机：%v , 邮箱：%v , 密码过期时间：%v .", user.Username, user.Nickname, user.Mobile, user.Email, utils.FormatTime(expireTime))
				switch config.Setting.Channel.ExpiredChannel {
				case "mail":
					subject := fmt.Sprintf("%s-密码过期通知", config.Setting.Channel.Mail.From)
					content := fmt.Sprintf(`<div>
                        <div>
                            %s, 您好!
                        </div>
                        <div style="padding: 8px 40px 8px 50px;">
                            <p>您的域控帐号(%s)的密码将于 %s 到期, 请登陆密码自助平台 %s 修改密码, 避免到期后导致无法使用。</p>
                            <p>请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
                        </div>
                        <div>
                            <p>此邮箱为系统邮箱，请勿回复。</p>
                        </div>
                    </div>`, user.Nickname, user.Username, utils.FormatTime(expireTime), config.Setting.Channel.PlatformUrl)

					if err := mail.SendMail([]string{user.Email}, subject, content); err != nil {
						logging.Logger().Sugar().Errorf("用户 %s 的密码过期通知使用 %s 方式发送失败，错误信息：%v", user.Username, config.Setting.Channel.ExpiredChannel, err)
					}
				case "aliyunsms":
					if err := aliyun.SendAliyunSMS(utils.FormatTime(expireTime), user.Mobile, config.Setting.Channel.AliyunSms.TemplateCodeExpired); err != nil {
						logging.Logger().Sugar().Errorf("用户 %s 的密码过期通知使用 %s 方式发送失败，错误信息：%v", user.Username, config.Setting.Channel.ExpiredChannel, err)
					}
				case "tencentsms":
					if err := tencent.SendTencentSMS(utils.FormatTime(expireTime), user.Mobile, config.Setting.Channel.TencentSms.TemplateCodeExpired); err != nil {
						logging.Logger().Sugar().Errorf("用户 %s 的密码过期通知使用 %s 方式发送失败，错误信息：%v", user.Username, config.Setting.Channel.ExpiredChannel, err)
					}
				default:
					return errors.New("不支持的消息通知通道")
				}

			}
		}

	}

	return nil
}
