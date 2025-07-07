package captcha

import (
	"fmt"
	"net/http"
	"password-self-service/internal/service/captcha"
	"password-self-service/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

var Handler handler

type handler struct{}

// sendCaptchaRequest 获取验证码
type sendCaptchaRequest struct {
	Username string `json:"username" bind:"username" binding:"required"`                        // 用户名
	Category string `json:"category" bind:"category" binding:"required,oneof=account password"` // 类别
}

// SendCaptcha
//
//	@Summary	获取验证码
//	@Tags		Captcha
//	@Accept		application/json
//	@Produce	json
//	@Security	BearerToken
//	@Param		sendCaptcha	body		sendCaptchaRequest	true	"获取验证码"
//	@Success	200			{object}	response.JsonRes
//	@Failure	500			{object}	response.JsonRes
//	@Router		/captcha/send [post]
func (h *handler) SendCaptcha(c *gin.Context) {
	var req sendCaptchaRequest

	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("%v:%v", "请求参数绑定错误", err)})
		return
	}

	err := captcha.Service.SendCaptcha(strings.ToLower(req.Username), req.Category)
	if err != nil {
		response.JsonExit(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Json(c, http.StatusOK, "获取验证码成功", nil)
}

// verifyCaptchaRequest 校验验证码结构体
type verifyCaptchaRequest struct {
	Username string `json:"username" bind:"username" binding:"required"`                        // 用户名
	Category string `json:"category" bind:"category" binding:"required,oneof=account password"` // 类别
	Code     string `json:"code" bind:"code" binding:"required"`                                // 验证码
}

// VerifyCaptcha
//
//	@Summary	校验验证码
//	@Tags		Captcha
//	@Accept		application/json
//	@Produce	json
//	@Security	BearerToken
//	@Param		VerifyCaptcha	body		verifyCaptchaRequest	true	"校验验证码"
//	@Success	200				{object}	response.JsonRes
//	@Failure	500				{object}	response.JsonRes
//	@Router		/captcha/verify [post]
func (h *handler) VerifyCaptcha(c *gin.Context) {
	var req verifyCaptchaRequest

	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JsonExit(c, http.StatusBadRequest, fmt.Sprintf("%v:%v", "请求参数绑定错误", err), nil)
		return
	}

	data, err := captcha.Service.VerifyCaptcha(strings.ToLower(req.Username), req.Category, req.Code)
	if err != nil {
		response.JsonExit(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Json(c, http.StatusOK, "验证码校验成功", data)
}
