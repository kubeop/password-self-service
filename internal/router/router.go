package router

import (
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"password-self-service/api/v1/account"
	"password-self-service/api/v1/captcha"
	"password-self-service/api/v1/password"
	"password-self-service/docs"
	"password-self-service/internal/middleware"
	"password-self-service/pkg/config"
	"password-self-service/views"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	//设置模式
	gin.SetMode(config.Setting.Server.Mode)

	// 创建带有默认中间件的路由:
	// 日志与恢复中间件
	server := gin.Default()

	// 启用全局跨域中间件
	server.Use(middleware.AllowCors())

	// 创建限速器
	limiter := tollbooth.NewLimiter(config.Setting.Server.RateLimit, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Second,
	})

	// 使用Gin限速中间件
	server.Use(middleware.LimitHandler(limiter))

	// 静态页面
	server.NoRoute(gin.WrapH(http.FileServer(http.FS(views.StaticAsset))))

	// 路由分组
	apiGroup := server.Group("/")

	InitAPIV0(apiGroup)
	InitAPIV1(apiGroup)

	return server
}

// InitAPIV0 基础接口
func InitAPIV0(r *gin.RouterGroup) {
	if config.Setting.Server.Mode == "debug" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.GET("health", func(ctx *gin.Context) { ctx.String(http.StatusOK, "ok") })
}

// InitAPIV1 业务接口
func InitAPIV1(r *gin.RouterGroup) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	router := r.Group("/api/v1")
	router.POST("/captcha/send", captcha.Handler.SendCaptcha)
	router.POST("/captcha/verify", captcha.Handler.VerifyCaptcha)
	router.POST("/unlock-account", account.Handler.UnlockAccount)
	router.POST("/reset-password", password.Handler.ResetPassword)
}
