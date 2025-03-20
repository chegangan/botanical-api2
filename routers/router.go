package routers

import (
	"botanical-api2/middleware"
	v1 "botanical-api2/routers/api/v1"
	"botanical-api2/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置所有路由
func SetupRouter(userService *service.UserService) *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	// 设置用户服务实例
	v1.SetUserService(userService)

	// 无需认证的API路由
	apiV1 := r.Group("/api/v1")
	{
		// 认证相关路由
		apiV1.POST("/auth/register", v1.RegisterUser)
		apiV1.POST("/auth/login", v1.LoginUser)

		// 需要JWT认证的路由
		authorized := apiV1.Group("/")
		authorized.Use(middleware.JWTMiddleware(userService))
		{
			// 用户资源路由
			authorized.GET("/users/:id", v1.GetUser)
			authorized.PUT("/users/:id", v1.UpdateUser)
			authorized.DELETE("/users/:id", v1.DeleteUser)
		}
	}

	return r
}
