package dto

// LoginRequest 登录请求DTO
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求DTO
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=20"` // 添加max=20限制
}

// AuthResponse 认证响应DTO
type AuthResponse struct {
	Token string      `json:"token"`
	User  UserSummary `json:"user"`
}
