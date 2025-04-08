package handlers

import (
	"botanical-api2/internal/api/middleware"
	"botanical-api2/internal/service"
	"botanical-api2/pkg/app"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 更新 Handler 结构体，添加新服务
type Handler struct {
	UserService        *service.UserService
	PictureService     service.PictureService
	FeedbackService    *service.FeedbackService
	PlantService       service.PlantVegetationService
	PlantClassService  service.PlantClassService
	ParkService        service.ParkService
	PlantNoticeService service.PlantNoticeService
}

// NewHandler 创建新的处理器
func NewHandler(service *service.Services) *Handler {
	return &Handler{
		UserService:        service.User,
		PictureService:     service.Picture,
		FeedbackService:    service.Feedback,
		PlantService:       service.PlantVegetation,
		PlantClassService:  service.PlantClass,
		ParkService:        service.Park,
		PlantNoticeService: service.PlantNotice,
	}
}

// RegisterRoutes 注册所有路由
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// 添加Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 静态文件服务 - 用于图片访问
	r.StaticFS("/uploads", http.Dir("./uploads"))

	// API版本分组
	apiV1 := r.Group("/api/v1")
	{
		// 健康检查接口
		apiV1.GET("/health", h.HealthCheck)

		// 注册公开访问的路由
		h.registerPublicRoutes(apiV1)

		// 注册需要JWT认证的路由
		h.registerAuthorizedRoutes(apiV1)
	}
}

// registerPublicRoutes 注册所有公开访问的路由
func (h *Handler) registerPublicRoutes(router *gin.RouterGroup) {
	// ===== 认证相关路由 =====
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.RegisterUser)
		auth.POST("/login", h.LoginUser)
	}

	// ===== 公开访问的图片路由 =====
	publicPictures := router.Group("/pictures")
	{
		publicPictures.GET("/:id", h.GetPicture)
	}

	// ===== 公开访问的用户资源 =====
	publicUsers := router.Group("/users")
	{
		publicUsers.GET("/:id/avatar", h.GetUserAvatar)
		publicUsers.GET("/:id/pictures", h.GetUserPictures)
	}

	// 植物公开路由
	plants := router.Group("/plants")
	{
		plants.GET("", h.GetPlants)
		plants.GET("/:id", h.GetPlant)
	}

	// 植物分类公开路由
	plantClasses := router.Group("/plant-classes")
	{
		plantClasses.GET("", h.GetPlantClasses)
		plantClasses.GET("/:id", h.GetPlantClass)
	}

	// 园区公开路由
	parks := router.Group("/parks")
	{
		parks.GET("", h.GetParks)
		parks.GET("/:id", h.GetPark)
		parks.GET("/:id/plants", h.GetParkPlants) // 使用相同的参数名
	}

	// 植物资讯公开路由
	notices := router.Group("/notices")
	{
		notices.GET("", h.GetPlantNotices)
		notices.GET("/:id", h.GetPlantNotice)
	}
}

// registerAuthorizedRoutes 注册所有需要认证的路由
func (h *Handler) registerAuthorizedRoutes(router *gin.RouterGroup) {
	// ===== 需要JWT认证的路由 =====
	authorized := router.Group("/")
	authorized.Use(middleware.JWTMiddleware(h.UserService))

	// 注册用户相关路由
	h.registerUserRoutes(authorized)

	// 注册资源管理路由
	h.registerResourceRoutes(authorized)

	// 注册管理员路由
	h.registerAdminRoutes(authorized)
}

// registerUserRoutes 注册用户相关路由
func (h *Handler) registerUserRoutes(router *gin.RouterGroup) {
	// 当前用户资源 (简化常用操作)
	me := router.Group("/me")
	{
		me.GET("", h.GetCurrentUser)                 // 获取当前用户信息
		me.PUT("", h.UpdateCurrentUser)              // 更新当前用户信息
		me.PUT("/password", h.ChangeCurrentPassword) // 修改当前用户密码
		me.POST("/avatar", h.UploadAvatar)           // 上传当前用户头像
		me.POST("/pictures", h.UploadUserPicture)    // 上传当前用户图片
		me.POST("/feedback", h.CreateFeedback)       // 提交用户反馈
	}

	// 用户资源管理
	users := router.Group("/users")
	{
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
		users.PUT("/:id/password", h.ChangePassword)
	}
}

// registerResourceRoutes 注册资源管理路由
func (h *Handler) registerResourceRoutes(router *gin.RouterGroup) {
	// 图片资源管理
	pictures := router.Group("/pictures")
	{
		pictures.DELETE("/:id", h.DeletePicture)
	}

	// 反馈资源管理
	feedback := router.Group("/feedback")
	{
		feedback.GET("/:id", h.GetFeedback)
		feedback.POST("", h.CreateFeedback)
		feedback.DELETE("/:id", h.DeleteFeedback)
		feedback.PUT("/:id/status", h.UpdateFeedbackStatus)
	}
}

// registerAdminRoutes 注册管理员路由
func (h *Handler) registerAdminRoutes(router *gin.RouterGroup) {
	// 为所有管理员路由添加管理员权限验证中间件
	admin := router.Group("")
	admin.Use(middleware.AdminRequired())
	{
		// 反馈管理
		feedbacks := admin.Group("/admin")
		{
			feedbacks.GET("/feedbacks", h.GetAllFeedbacks)
		}

		// 植物管理
		plants := admin.Group("/plants")
		{
			plants.POST("", h.CreatePlant)
			plants.PUT("/:id", h.UpdatePlant)
			plants.DELETE("/:id", h.DeletePlant)
			plants.POST("/:id/images", h.UploadPlantImage)
		}

		// 植物分类管理
		plantClasses := admin.Group("/plant-classes")
		{
			plantClasses.POST("", h.CreatePlantClass)
			plantClasses.PUT("/:id", h.UpdatePlantClass)
			plantClasses.DELETE("/:id", h.DeletePlantClass)
			plantClasses.POST("/:id/images", h.UploadPlantClassImage)
		}

		// 管理员路由部分
		adminParks := admin.Group("/parks")
		{
			adminParks.POST("", h.CreatePark)
			adminParks.PUT("/:id", h.UpdatePark)
			adminParks.DELETE("/:id", h.DeletePark)
			adminParks.POST("/:id/images", h.UploadParkImage)
			adminParks.POST("/:id/plants/:plant_id", h.AddPlantToPark)        // 使用相同的参数名:id
			adminParks.DELETE("/:id/plants/:plant_id", h.RemovePlantFromPark) // 使用相同的参数名:id
		}

		// 植物资讯管理
		notices := admin.Group("/notices")
		{
			notices.POST("", h.CreatePlantNotice)
			notices.PUT("/:id", h.UpdatePlantNotice)
			notices.DELETE("/:id", h.DeletePlantNotice)
			notices.POST("/:id/images", h.UploadPlantNoticeImage)
			notices.PUT("/:id/recommend", h.TogglePlantNoticeRecommend)
		}
	}
}

// HealthCheck 健康检查
// @Summary API健康状态
// @Description 获取API服务的健康状态信息
// @Tags 健康检查
// @Accept json
// @Produce json
// @SUCCESS 200 {object} app.Response{data=map[string]interface{}} "API服务正常"
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	info := map[string]interface{}{
		"status":      "running",
		"timestamp":   time.Now().Format(time.RFC3339),
		"go_version":  runtime.Version(),
		"go_os":       runtime.GOOS,
		"go_arch":     runtime.GOARCH,
		"cpu_num":     runtime.NumCPU(),
		"goroutines":  runtime.NumGoroutine(),
		"api_version": "v1",
	}

	app.SUCCESS(c, info)
}
