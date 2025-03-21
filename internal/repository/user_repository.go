package repository

import (
	"botanical-api2/internal/models"

	"github.com/jinzhu/gorm"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *gorm.DB // 改为小写，使其成为私有字段
}

// NewUserRepository 创建新的用户仓库
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db} // 使用小写字段名
}

// CreateUser 创建用户
func (repo *UserRepository) CreateUser(user *models.User) error {
	return repo.db.Create(user).Error // 使用小写字段
}

// GetUserByID 根据ID获取用户
func (repo *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := repo.db.Where("id = ?", id).First(&user).Error // 使用小写字段
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserFields 动态更新用户指定字段
func (repo *UserRepository) UpdateUserFields(id int, fields map[string]interface{}) error {
	return repo.db.Model(&models.User{}).Where("id = ?", id).Updates(fields).Error // 使用小写字段
}

// DeleteUser 删除用户
func (repo *UserRepository) DeleteUser(id int) error {
	return repo.db.Where("id = ?", id).Delete(&models.User{}).Error // 使用小写字段
}

// CheckPhoneExists 检查手机号是否已存在
func (repo *UserRepository) CheckPhoneExists(phone string) (bool, error) {
	var count int
	err := repo.db.Model(&models.User{}). // 使用小写字段
						Where("phone <> '' AND phone = ?", phone).
						Count(&count).Error

	return count > 0, err
}

// GetUserByPhone 根据手机号获取用户
func (repo *UserRepository) GetUserByPhone(phone string) (*models.User, error) {
	var user models.User
	err := repo.db.Where("phone = ?", phone).First(&user).Error // 使用小写字段
	if err != nil {
		return nil, err
	}
	return &user, nil
}
