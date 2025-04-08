package service

import "botanical-api2/internal/repository"

// Services 服务集合，管理所有业务服务
type Services struct {
	User            *UserService
	Picture         PictureService
	Feedback        *FeedbackService
	PlantVegetation PlantVegetationService
	PlantClass      PlantClassService
	Park            ParkService
	PlantNotice     PlantNoticeService
	// 将来可以添加其他服务
}

// NewServices 创建服务集合实例
func NewServices(repository *repository.Repositories) *Services {
	return &Services{
		User:            NewUserService(repository.User),
		Picture:         NewPictureService(repository),
		Feedback:        NewFeedbackService(repository),
		PlantVegetation: NewPlantVegetationService(repository),
		PlantClass:      NewPlantClassService(repository),
		Park:            NewParkService(repository),
		PlantNotice:     NewPlantNoticeService(repository),
	}
}
