package repository

import (
	"botanical-api2/models" // Update with your actual project path

	"github.com/jinzhu/gorm"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository 创建新的用户仓库
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser 创建用户
func (repo *UserRepository) CreateUser(user *models.User) error {
	return repo.DB.Create(user).Error
}

// GetUserByID 根据ID获取用户
func (repo *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := repo.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (repo *UserRepository) UpdateUser(user *models.User) error {
	return repo.DB.Save(user).Error
}

// DeleteUser 删除用户
func (repo *UserRepository) DeleteUser(id int) error {
	return repo.DB.Where("id = ?", id).Delete(&models.User{}).Error
}

// CheckUserExist 检查用户名或手机号是否已存在
func (repo *UserRepository) CheckUserExist(username, phone string) (bool, error) {
	var count int
	err := repo.DB.Model(&models.User{}).
		Where("username = ? OR (phone <> '' AND phone = ?)", username, phone).
		Count(&count).Error

	return count > 0, err
}

// GetUserByUsername 根据用户名获取用户
func (repo *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := repo.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
