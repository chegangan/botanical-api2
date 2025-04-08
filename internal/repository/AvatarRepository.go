package repository

import (
	"errors"

	"botanical-api2/internal/models"

	"github.com/jinzhu/gorm"
)

// AvatarRepository 头像存储库实现
type AvatarRepository struct {
	DB *gorm.DB
}

// NewAvatarRepository 创建头像存储库
func NewAvatarRepository(db *gorm.DB) *AvatarRepository {
	return &AvatarRepository{DB: db}
}

// Create 创建头像
func (r *AvatarRepository) Create(avatar *models.UserAvatar) error {
	return r.DB.Create(avatar).Error
}

// GetByID 根据ID获取头像
func (r *AvatarRepository) GetByID(id int) (*models.UserAvatar, error) {
	var avatar models.UserAvatar
	if err := r.DB.Where("id = ?", id).First(&avatar).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &avatar, nil
}

// GetByUserID 根据用户ID获取头像
func (r *AvatarRepository) GetByUserID(userID int) (*models.UserAvatar, error) {
	var avatar models.UserAvatar
	if err := r.DB.Where("user_id = ?", userID).First(&avatar).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &avatar, nil
}

// Update 更新头像
func (r *AvatarRepository) Update(avatar *models.UserAvatar) error {
	return r.DB.Save(avatar).Error
}

// Delete 删除头像
func (r *AvatarRepository) Delete(id int) error {
	return r.DB.Delete(&models.UserAvatar{}, id).Error
}
