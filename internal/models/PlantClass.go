package models

// PlantClass 植物分类模型
// @Description 植物分类信息
type PlantClass struct {
	ID      int    `json:"id" gorm:"primary_key;auto_increment;unique"` // 主键ID，自增长
	SrcList string `json:"src_list" gorm:"column:src_list;size:255"`    // 分类相关图片列表(多图URL，以逗号分隔)
	Name    string `json:"name" gorm:"column:name;size:255;not null"`   // 分类名称
}

// TableName 设置表名
func (PlantClass) TableName() string {
	return "plant_class"
}
