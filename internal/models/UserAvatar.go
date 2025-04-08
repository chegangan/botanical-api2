package models

// UserAvatar 用户头像模型
// @Description 用户头像信息
type UserAvatar struct {
	ID     int    `json:"id" gorm:"primary_key;auto_increment;unique"` // 头像ID
	URL    string `json:"url" gorm:"column:url;size:255"`              // 头像URL
	UserID int    `json:"user_id" gorm:"column:user_id;index"`         // 用户ID，外键
}

// TableName 设置表名
func (UserAvatar) TableName() string {
	return "user_avatar"
}
