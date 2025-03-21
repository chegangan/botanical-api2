package api // 修改包名为api，与目录结构一致

import (
	"botanical-api2/internal/api/handlers"
	"botanical-api2/internal/api/middleware"
	"botanical-api2/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupHandler 设置所有路由
func SetupHandler(userService *service.UserService) *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	// 设置用户服务实例
	h := handlers.NewHandler(userService)

	// 无需认证的API路由
	apiV1 := r.Group("/api/v1")
	{
		// 认证相关路由
		apiV1.POST("/auth/register", h.RegisterUser)
		apiV1.POST("/auth/login", h.LoginUser)

		// 需要JWT认证的路由
		authorized := apiV1.Group("/")
		authorized.Use(middleware.JWTMiddleware(userService)) // 直接传递userService
		{
			// 用户资源路由
			authorized.GET("/users/:id", h.GetUser)
			authorized.PUT("/users/:id", h.UpdateUser)
			authorized.DELETE("/users/:id", h.DeleteUser)
			authorized.PUT("/users/:id/password", h.ChangePassword)
		}
	}

	return r
}
