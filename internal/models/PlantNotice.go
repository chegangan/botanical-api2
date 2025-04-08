package models

import "time"

// PlantNotice 植物资讯模型
// @Description 植物资讯信息
type PlantNotice struct {
	ID          int       `json:"id" gorm:"primary_key;auto_increment;unique"` // 主键ID，自增长
	Intro       string    `json:"intro" gorm:"column:intro;type:text"`         // 资讯内容
	Src         string    `json:"src" gorm:"column:src;size:255"`              // 资讯图片URL
	SkimSort    int       `json:"skim_sort" gorm:"column:skim_sort"`           // 浏览量
	IsRecommend int       `json:"is_recommend" gorm:"column:is_recommend"`     // 是否推荐 0不推荐 1推荐
	Title       string    `json:"title" gorm:"column:title;size:255"`          // 资讯标题
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`         // 创建时间
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`         // 更新时间
}

// TableName 设置表名
func (PlantNotice) TableName() string {
	return "plant_notice"
}
