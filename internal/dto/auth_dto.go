package dto

// LoginRequest 登录请求DTO
// @Description 用户登录的请求参数
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required" example:"13800138000"`    // 用户手机号
	Password string `json:"password" binding:"required" example:"Password123"` // 用户密码
}

// RegisterRequest 注册请求DTO
// @Description 用户注册的请求参数
type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"zhangsan"`                 // 用户名称
	Phone    string `json:"phone" binding:"required" example:"13800138000"`                 // 用户手机号
	Password string `json:"password" binding:"required,min=6,max=20" example:"Password123"` // 用户密码
}

// AuthResponse 认证响应DTO
// @Description 认证成功的响应数据
type AuthResponse struct {
	Token string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT令牌
	User  UserSummary `json:"user"`                                                    // 用户信息摘要
}
