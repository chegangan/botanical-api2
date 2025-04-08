package repository

import (
	"errors"

	"botanical-api2/internal/models"

	"github.com/jinzhu/gorm"
)

// PictureRepository 图片存储库实现
type PictureRepository struct {
	DB *gorm.DB
}

// NewPictureRepository 创建图片存储库
func NewPictureRepository(db *gorm.DB) *PictureRepository {
	return &PictureRepository{DB: db}
}

// Create 创建图片
func (r *PictureRepository) Create(picture *models.UserPicture) error {
	return r.DB.Create(picture).Error
}

// GetByID 根据ID获取图片
func (r *PictureRepository) GetByID(id int) (*models.UserPicture, error) {
	var picture models.UserPicture
	if err := r.DB.Where("id = ?", id).First(&picture).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &picture, nil
}

// GetByUserID 根据用户ID获取所有图片
func (r *PictureRepository) GetByUserID(userID int) ([]models.UserPicture, error) {
	var pictures []models.UserPicture
	if err := r.DB.Where("user_id = ?", userID).Find(&pictures).Error; err != nil {
		return nil, err
	}
	return pictures, nil
}

// Delete 删除图片
func (r *PictureRepository) Delete(id int) error {
	return r.DB.Delete(&models.UserPicture{}, id).Error
}
