package models

import (
	"time"
)

// User 用户模型
// @Description 用户数据库模型
type User struct {
	ID              int        `json:"id" gorm:"primary_key" example:"1"`                                  // 用户ID
	Username        string     `json:"username" gorm:"size:50;not null;default:'user'" example:"zhangsan"` // 用户名称
	Phone           string     `json:"phone" gorm:"size:20;unique;not null" example:"13800138000"`         // 手机号码
	PasswordHash    string     `json:"-" gorm:"column:password_hash;size:255;not null"`                    // 密码哈希(不返回给前端)
	UserRole        int8       `json:"user_role" gorm:"default:1" example:"1"`                             // 用户角色(1-普通用户,2-VIP会员,9-管理员)
	UserStatus      int8       `json:"user_status" gorm:"default:1"`                                       // 用户状态(0-禁用,1-正常,2-待验证)
	Token           string     `json:"-" gorm:"size:255"`                                                  // JWT令牌(不返回给前端)
	TokenExpireTime *time.Time `json:"-" gorm:"column:token_expire_time"`                                  // 令牌过期时间(不返回给前端)
	CreatedAt       time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`                        // 创建时间
	UpdatedAt       time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`                        // 更新时间

	// 关联
	Avatar *UserAvatar `json:"avatar,omitempty" gorm:"foreignKey:UserID"` // 用户头像
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}
