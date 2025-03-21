package dto

// UserSummary 用户信息摘要DTO
// @Description 用户基本信息，不包含敏感字段
type UserSummary struct {
	ID       int    `json:"id" example:"1"`              // 用户ID
	Username string `json:"username" example:"zhangsan"` // 用户名称
	Phone    string `json:"phone" example:"13800138000"` // 用户手机号
	UserRole int8   `json:"user_role" example:"1"`       // 用户角色(1:普通用户,9:管理员)
}

// PasswordChangeRequest 密码修改请求DTO
// @Description 修改密码的请求参数
type PasswordChangeRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=20" example:"OldPass123"` // 原密码
	NewPassword string `json:"new_password" binding:"required,min=6,max=20" example:"NewPass456"` // 新密码
}
