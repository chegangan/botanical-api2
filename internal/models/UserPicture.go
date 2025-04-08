package models

// UserPicture 用户图片模型
// @Description 用户上传的照片
type UserPicture struct {
	ID     int    `json:"id" gorm:"primary_key;auto_increment;unique"` // 图片ID
	URL    string `json:"url" gorm:"column:url;size:255"`              // 图片URL
	UserID int    `json:"user_id" gorm:"column:user_id;index"`         // 用户ID，外键
}

// TableName 设置表名
func (UserPicture) TableName() string {
	return "user_picture"
}
