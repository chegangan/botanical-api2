package main

import (
	"botanical-api2/internal/api"
	"botanical-api2/internal/models"
	"botanical-api2/internal/repository"
	"botanical-api2/internal/service"
	"botanical-api2/pkg/setting"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jinzhu/gorm"
)

// Server 应用服务器结构体
type Server struct {
	db       *gorm.DB
	router   *http.Server
	repos    *repository.Repositories
	services *service.Services
}

func NewServer() *Server {
	// 初始化数据库
	models.InitDB()
	db := models.GetDB()

	// 创建仓库集合
	repos := repository.NewRepositories(db)

	// 创建服务
	userService := service.NewUserService(repos.User)
	services := service.NewServices(userService)

	// 设置路由
	Handler := api.SetupHandler(services.User)

	// 创建HTTP服务器
	httpServer := &http.Server{
		Addr:           ":" + strconv.Itoa(setting.HttpPort),
		Handler:        Handler,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: setting.MaxHeaderBytes,
	}

	return &Server{
		db:       db,
		router:   httpServer,
		repos:    repos, // 保存仓库集合实例
		services: services,
	}
}

// Run 启动服务器
func (s *Server) Run() error {
	// 创建通道监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 在后台运行服务器
	go func() {
		log.Printf("服务器运行在 %s 模式，端口 %d", setting.RunMode, setting.HttpPort)
		if err := s.router.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("监听错误: %s\n", err)
		}
	}()

	// 等待中断信号
	<-quit
	log.Println("正在关闭服务器...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := s.router.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}

	// 关闭数据库连接
	if err := s.db.Close(); err != nil {
		log.Fatal("数据库关闭错误:", err)
	}

	log.Println("服务器优雅关闭")
	return nil
}
