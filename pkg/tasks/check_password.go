package tasks

import (
	"context"
	"password-self-service/internal/service/password"
	"password-self-service/pkg/config"
	"password-self-service/pkg/logging"

	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	ctx = gctx.New()
)

func Init() {
	cron := gcron.New()

	if config.Setting.Cron.Enabled {
		CheckPasswordExpired()
		cron.Start("CheckPasswordExpired")
	}

}

func CheckPasswordExpired() {
	_, err := gcron.Add(ctx, config.Setting.Cron.Schedule, func(ctx context.Context) {
		err := password.Service.PasswordExpired()
		if err != nil {
			logging.Logger().Sugar().Errorf("密码过期检查失败:%v", err)
		}

	}, "CheckPasswordExpired")
	if err != nil {
		logging.Logger().Sugar().Errorf("定时任务添加失败:%v", err)
	}

}
