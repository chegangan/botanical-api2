package dto

// UserSummary 用户信息摘要DTO（用于敏感场景返回）
type UserSummary struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	UserRole int8   `json:"user_role"`
}

// PasswordChangeRequest 密码修改请求DTO
type PasswordChangeRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=20"` // 添加max=20限制
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"` // 添加max=20限制
}
