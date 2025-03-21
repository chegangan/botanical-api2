package handlers

import "botanical-api2/internal/service"

// Handler 处理器结构体，包含所有服务依赖
type Handler struct {
	UserService *service.UserService
}

// NewHandler 创建新的处理器
func NewHandler(userService *service.UserService) *Handler {
	return &Handler{
		UserService: userService,
	}
}
