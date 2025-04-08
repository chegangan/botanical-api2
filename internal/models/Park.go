package models

// Park 园区模型
// @Description 园区信息模型
type Park struct {
	ID  int    `json:"id" gorm:"primary_key;auto_increment;unique"` // 园区图片ID
	Src string `json:"src" gorm:"column:src;size:255;not null"`     // 园区区域范围图片链接或路径
}

// TableName 设置表名
func (Park) TableName() string {
	return "park"
}
