package models

// UserCollect 用户收藏模型
// @Description 用户收藏信息
type UserCollect struct {
	ID        int `json:"id" gorm:"primary_key;auto_increment;unique"` // 收藏ID
	UserID    int `json:"user_id" gorm:"column:user_id;index"`         // 用户ID，外键
	Type      int `json:"type" gorm:"column:type"`                     // 收藏类型
	CollectID int `json:"collect_id" gorm:"column:collect_id"`         // 收藏的对象ID
}

// TableName 设置表名
func (UserCollect) TableName() string {
	return "user_collect"
}
