package middleware

import (
	"botanical-api2/internal/models"
	"botanical-api2/pkg/app"
	"botanical-api2/pkg/e"

	"github.com/gin-gonic/gin"
)

// AdminRequired 检查用户是否具有管理员权限的中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取已认证的用户信息
		// JWT中间件应该已经将用户信息设置到上下文中
		userInterface, exists := c.Get("user")
		if !exists {
			app.Error(c, e.UNAUTHORIZED, e.GetMsg(e.UNAUTHORIZED))
			c.Abort()
			return
		}

		// 类型断言获取用户对象
		user, ok := userInterface.(*models.User)
		if !ok {
			app.Error(c, e.BAD_REQUEST, "用户信息类型错误")
			c.Abort()
			return
		}

		// 检查用户是否具有管理员角色
		if user.UserRole != 9 { // 假设9是管理员角色值
			app.Error(c, e.FORBIDDEN, e.GetMsg(e.FORBIDDEN))
			c.Abort()
			return
		}

		// 用户是管理员，继续处理请求
		c.Next()
	}
}
