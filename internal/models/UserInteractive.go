package models

// UserInteractive 用户互动模型
// @Description 用户互动信息
type UserInteractive struct {
	ID      int    `json:"id" gorm:"primary_key;auto_increment;unique"` // 互动ID
	UserID  int    `json:"user_id" gorm:"column:user_id;index"`         // 用户ID，外键
	Intro   string `json:"intro" gorm:"column:intro;type:text"`         // 介绍内容
	SrcList string `json:"src_list" gorm:"column:src_list;size:2000"`   // 图片列表，逗号分隔的URL
	Title   string `json:"title" gorm:"column:title;size:255"`          // 标题
}

// TableName 设置表名
func (UserInteractive) TableName() string {
	return "user_interactive"
}
