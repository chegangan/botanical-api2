package models

// PlantParkVegetation 园区与植物的关联模型
// @Description 记录园区与植物之间的多对多关系
type PlantParkVegetation struct {
	ID      int `json:"id" gorm:"primary_key;auto_increment;unique"` // 主键ID
	ParkID  int `json:"park_id" gorm:"column:park_id;not null"`      // 园区ID
	PlantID int `json:"plant_id" gorm:"column:plant_id;not null"`    // 植物ID
}

// TableName 设置表名
func (PlantParkVegetation) TableName() string {
	return "plant_park_vegetation"
}
