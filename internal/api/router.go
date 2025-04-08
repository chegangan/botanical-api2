package api

import (
	"botanical-api2/internal/api/handlers"
	"botanical-api2/internal/api/middleware"
	"botanical-api2/internal/service"

	_ "botanical-api2/docs"

	"github.com/gin-gonic/gin"
)

/**
使用swagger的命令
swag init --parseDependency --parseInternal -g router.go --output ../../docs
http://localhost:8000/swagger/index.html
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
func SetupHandler(service *service.Services) *gin.Engine {
	// 创建Gin引擎
	r := gin.New()

	// 添加全局中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	// 创建handler实例
	handler := handlers.NewHandler(service)

	// 注册所有路由
	handler.RegisterRoutes(r)

	return r
}
