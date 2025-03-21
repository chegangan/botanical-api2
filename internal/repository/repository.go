package repository

import "github.com/jinzhu/gorm"

// Repositories 仓库集合，统一管理所有数据访问层实例
type Repositories struct {
	User *UserRepository
	// 将来可添加的其他仓库
	// Product  *ProductRepository
	// Order    *OrderRepository
	// Category *CategoryRepository
}

// NewRepositories 创建仓库集合实例
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
	}
}
