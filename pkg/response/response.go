package response

import "github.com/gin-gonic/gin"

// JsonRes 数据返回通用JSON数据结构
type JsonRes struct {
	Code    int         `json:"code"`    // 错误码((0:成功, 1:失败, >1:错误码))
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 返回数据(业务接口定义具体数据结构)
}

// Json 返回标准JSON数据
func Json(c *gin.Context, httpCode int, message string, data interface{}) {
	c.JSON(httpCode, JsonRes{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// JsonExit 返回标准JSON数据并退出当前HTTP执行函数。
func JsonExit(c *gin.Context, httpCode int, message string, data interface{}) {
	c.JSON(httpCode, JsonRes{
		Code:    -1,
		Message: message,
		Data:    data,
	})
}
