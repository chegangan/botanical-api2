package models

import (
	"time"
)

// User 用户模型
// @Description 用户数据库模型
type User struct {
	ID              int        `json:"id" gorm:"primary_key" example:"1"`                          // 用户ID
	Username        string     `json:"username" gorm:"size:50;not null" example:"zhangsan"`        // 用户名称
	Phone           string     `json:"phone" gorm:"size:20;unique;not null" example:"13800138000"` // 手机号码
	PasswordHash    string     `json:"-" gorm:"size:100;not null"`                                 // 密码哈希(不返回给前端)
	UserRole        int8       `json:"user_role" gorm:"default:1" example:"1"`                     // 用户角色
	Token           string     `json:"-" gorm:"size:500"`                                          // JWT令牌(不返回给前端)
	TokenExpireTime *time.Time `json:"-"`                                                          // 令牌过期时间(不返回给前端)
	CreatedAt       time.Time  `json:"created_at"`                                                 // 创建时间
	UpdatedAt       time.Time  `json:"updated_at"`                                                 // 更新时间
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}
