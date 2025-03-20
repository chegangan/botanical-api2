package main

import (
	"botanical-api2/models"
	"botanical-api2/pkg/setting"
	"botanical-api2/routers"
	"botanical-api2/service"
	"log"
	"net/http"
	"strconv"
)

func main() {
	// 初始化数据库
	models.InitDB()
	db := models.GetDB()

	// 创建用户服务实例
	userService := service.NewUserService(db)

	// 设置路由
	r := routers.SetupRouter(userService)

	// 创建自定义服务器，使用配置文件中的所有参数
	server := &http.Server{
		Addr:           ":" + strconv.Itoa(setting.HttpPort),
		Handler:        r,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: setting.MaxHeaderBytes,
	}

	// 启动服务器
	log.Printf("服务器运行在 %s 模式，端口 %d", setting.RunMode, setting.HttpPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
