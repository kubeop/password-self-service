package password

import (
	"fmt"
	"net/http"
	"password-self-service/internal/service/password"
	"password-self-service/pkg/response"

	"github.com/gin-gonic/gin"
)

var Handler handler

type handler struct{}

// resetPasswordRequest 重置密码结构体
type resetPasswordRequest struct {
	Username    string `json:"username" bind:"username" binding:"required,min=2,max=20"` // 用户名
	NewPassword string `json:"newPassword" bind:"newPassword" binding:"required"`        // 新密码
	Code        string `json:"code" bind:"code" binding:"required"`                      // 验证码
}

// ResetPassword
//
//	@Summary	重置密码
//	@Tags		User
//	@Accept		application/json
//	@Produce	json
//	@Security	BearerToken
//	@Param		ResetPassword	body		resetPasswordRequest	true	"重置密码"
//	@Success	201				{object}	response.JsonRes
//	@Failure	500				{object}	response.JsonRes
//	@Router		/reset-password [post]
func (h *handler) ResetPassword(c *gin.Context) {
	var req resetPasswordRequest

	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JsonExit(c, http.StatusBadRequest, fmt.Sprintf("%v:%v", "请求参数绑定错误", err), nil)
		return
	}

	// // RSA 解密
	// rsa.ReadRSAKey()
	// plainText, err := rsa.Decrypt(req.NewPassword)
	// if err != nil {
	// 	response.JsonExit(c, http.StatusInternalServerError, err.Error(), nil)
	// 	return
	// }

	err := password.Service.ResetPassword(
		req.Username,
		req.NewPassword,
		req.Code,
	)
	if err != nil {
		response.JsonExit(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Json(c, http.StatusCreated, "重置密码成功", nil)
}
