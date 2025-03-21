package models

import (
	"time"
)

// User 用户表模型
type User struct {
	ID              int        `gorm:"primary_key;auto_increment;unique" json:"id"`
	Username        string     `gorm:"size:50;not null;default:'user'" json:"username"`
	Phone           string     `gorm:"size:20" json:"phone" binding:"omitempty"`
	PasswordHash    string     `gorm:"size:255;not null" json:"-"` // json中使用'-'忽略此字段，不返回给前端
	Token           string     `gorm:"size:255" json:"-"`
	TokenExpireTime *time.Time `gorm:"" json:"-"`
	UserStatus      int8       `gorm:"default:1" json:"user_status" example:"1"` // 0-禁用,1-正常,2-待验证
	UserRole        int8       `gorm:"default:1" json:"user_role" example:"1"`   // 1-普通用户,2-VIP会员,9-管理员
	CreatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}
