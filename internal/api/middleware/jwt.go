package middleware

import (
	"botanical-api2/internal/service"
	"botanical-api2/pkg/app" // 添加这一行
	"botanical-api2/pkg/e"   // 添加这一行
	"botanical-api2/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware JWT认证中间件
func JWTMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			app.Error(c, e.ERROR_USER_NO_PERMISSION, "缺失Authorization头信息")
			c.Abort()
			return
		}

		// 移除Bearer前缀
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// 使用jwt包解析令牌
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			app.Error(c, e.ERROR_USER_NO_PERMISSION, err.Error())
			c.Abort()
			return
		}

		// 获取用户信息
		user, err := userService.GetUserByID(claims.ID)
		if err != nil {
			app.Error(c, e.ERROR_USER_NOT_FOUND, "用户不存在")
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user", user)
		c.Set("claims", claims) // 可选：将claims也存入上下文
		c.Next()
	}
}
