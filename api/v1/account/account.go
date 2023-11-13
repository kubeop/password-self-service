package account

import (
	"fmt"
	"net/http"
	"password-self-service/internal/service/account"
	"password-self-service/pkg/response"

	"github.com/gin-gonic/gin"
)

var Handler handler

type handler struct{}

// unlockAccountRequest 解锁AD用户结构体
type unlockAccountRequest struct {
	Username string `json:"username" bind:"username" binding:"required,min=2,max=20"` // 用户名
	Code     string `json:"code" bind:"code" binding:"required"`                      // 验证码
}

// UnlockAccount
//
//	@Summary	解锁账户
//	@Tags		User
//	@Accept		application/json
//	@Produce	json
//	@Security	BearerToken
//	@Param		UnlockAccount	body		unlockAccountRequest	true	"解锁账户"
//	@Success	201				{object}	response.JsonRes
//	@Failure	500				{object}	response.JsonRes
//	@Router		/unlock-account [post]
func (h *handler) UnlockAccount(c *gin.Context) {
	var req unlockAccountRequest

	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JsonExit(c, http.StatusBadRequest, fmt.Sprintf("%v:%v", "请求参数绑定错误", err), nil)
		return
	}

	err := account.Service.UnlockAccount(
		req.Username,
		req.Code,
	)
	if err != nil {
		response.JsonExit(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Json(c, http.StatusCreated, "解锁AD用户成功", nil)
}
