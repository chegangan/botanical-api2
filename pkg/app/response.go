package app

import (
	"botanical-api2/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应格式
// @Description API统一响应格式
type Response struct {
	Code    int         `json:"code" example:"200"`     // 状态码
	Message string      `json:"message" example:"操作成功"` // 消息
	Data    interface{} `json:"data"`                   // 数据
}

// Success 响应成功
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    e.Success,
		Message: e.GetMsg(e.Success), // 使用统一的成功消息
		Data:    data,
	})
}

// Error 响应错误
func Error(c *gin.Context, code int, details string) {
	// 从错误码获取标准错误消息
	standardMsg := e.GetMsg(code)

	// 如果提供了详细信息且不同于标准消息，则将其附加到标准消息上
	message := standardMsg
	if details != "" && details != standardMsg {
		message = standardMsg + ": " + details
	}

	// 根据错误码确定HTTP状态码
	httpStatus := getHttpStatusByCode(code)

	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// 根据业务错误码获取适当的HTTP状态码
func getHttpStatusByCode(code int) int {
	switch code {
	case e.ErrorInvalidParams:
		return http.StatusBadRequest
	case e.ErrorUserNotFound:
		return http.StatusNotFound
	case e.ErrorUserUnauthorized:
		return http.StatusUnauthorized
	default:
		if code >= 10000 { // 自定义业务错误码
			return http.StatusBadRequest
		}
		return http.StatusInternalServerError
	}
}
