package models

// PlantVegetation 植物详细信息模型
// @Description 植物详细信息
type PlantVegetation struct {
	ID       int    `json:"id" gorm:"primary_key;auto_increment;unique"`       // 主键ID，自增长
	ClassID  int    `json:"class_id" gorm:"column:class_id;not null"`          // 关联植物分类ID
	Location string `json:"location" gorm:"column:location;size:255;not null"` // 植物所在位置/分布区域
	SrcList  string `json:"src_list" gorm:"column:src_list;size:255"`          // 植物相关图片列表(多图URL，以逗号分隔)
	Src      string `json:"src" gorm:"column:src;size:255;not null"`           // 植物主图URL
	Intro    string `json:"intro" gorm:"column:intro;type:text"`               // 植物详细介绍
	SkimSort int    `json:"skim_sort" gorm:"column:skim_sort"`                 // 浏览排序权重(数字越大排序越靠前)
	QrCode   string `json:"qr_code" gorm:"column:qr_code;size:255"`            // 植物专属二维码URL
	Name     string `json:"name" gorm:"column:name;size:255;not null"`         // 植物名称

	// 关联
	Class *PlantClass `json:"class,omitempty" gorm:"foreignKey:ClassID"` // 植物分类信息
}

// TableName 设置表名
func (PlantVegetation) TableName() string {
	return "plant_vegetation"
}
