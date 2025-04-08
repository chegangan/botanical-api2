package repository

import "github.com/jinzhu/gorm"

// Repositories 仓库集合，统一管理所有数据访问层实例
type Repositories struct {
	User            *UserRepository
	Avatar          *AvatarRepository
	Picture         *PictureRepository
	Feedback        *FeedbackRepository
	PlantVegetation *PlantVegetationRepository
	PlantClass      *PlantClassRepository
	Park            *ParkRepository
	PlantNotice     *PlantNoticeRepository
	// 将来可添加的其他仓库
}

// NewRepositories 创建仓库集合实例
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:            NewUserRepository(db),
		Avatar:          NewAvatarRepository(db),
		Picture:         NewPictureRepository(db),
		Feedback:        NewFeedbackRepository(db),
		PlantVegetation: NewPlantVegetationRepository(db),
		PlantClass:      NewPlantClassRepository(db),
		Park:            NewParkRepository(db),
		PlantNotice:     NewPlantNoticeRepository(db),
	}
}
