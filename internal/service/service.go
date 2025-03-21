package service

// Services 服务集合，管理所有业务服务
type Services struct {
	User *UserService
	// 将来可以添加其他服务
	// Auth      *AuthService
	// Email     *EmailService
	// Storage   *StorageService
}

// NewServices 创建服务集合实例
func NewServices(userService *UserService) *Services {
	return &Services{
		User: userService,
	}
}
