package api

import (
	"botanical-api2/internal/api/handlers"
	"botanical-api2/internal/api/middleware"
	"botanical-api2/internal/service"

	_ "botanical-api2/docs" // 这行很重要，导入自动生成的docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/**
使用swagger的命令
swag init --parseDependency --parseInternal -g router.go --output ../../docs
*/

// @title Botanical API
// @version 1.0
// @description 植物管理系统API服务
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email your-email@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8000
// @BasePath /api/v1
// @schemes http
func SetupHandler(userService *service.UserService) *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	// 设置用户服务实例
	h := handlers.NewHandler(userService)

	// 添加Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 无需认证的API路由
	apiV1 := r.Group("/api/v1")
	{
		// 认证相关路由
		apiV1.POST("/auth/register", h.RegisterUser)
		apiV1.POST("/auth/login", h.LoginUser)

		// 需要JWT认证的路由
		authorized := apiV1.Group("/")
		authorized.Use(middleware.JWTMiddleware(userService))
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
