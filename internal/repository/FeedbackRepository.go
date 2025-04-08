package repository

import (
	"errors"

	"botanical-api2/internal/models"

	"github.com/jinzhu/gorm"
)

// FeedbackRepository 用户反馈存储库实现
type FeedbackRepository struct {
	DB *gorm.DB
}

// NewFeedbackRepository 创建用户反馈存储库
func NewFeedbackRepository(db *gorm.DB) *FeedbackRepository {
	return &FeedbackRepository{DB: db}
}

// Create 创建用户反馈
func (r *FeedbackRepository) Create(feedback *models.UserFeedback) error {
	return r.DB.Create(feedback).Error
}

// GetByID 根据ID获取反馈
func (r *FeedbackRepository) GetByID(id int) (*models.UserFeedback, error) {
	var feedback models.UserFeedback
	if err := r.DB.Where("id = ?", id).First(&feedback).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &feedback, nil
}

// GetByUserID 根据用户ID获取所有反馈
func (r *FeedbackRepository) GetByUserID(userID int) ([]models.UserFeedback, error) {
	var feedbacks []models.UserFeedback
	if err := r.DB.Where("user_id = ?", userID).Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
}

// GetAll 获取所有反馈
func (r *FeedbackRepository) GetAll() ([]models.UserFeedback, error) {
	var feedbacks []models.UserFeedback
	if err := r.DB.Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
}

// Update 更新反馈信息
func (r *FeedbackRepository) Update(feedback *models.UserFeedback) error {
	return r.DB.Save(feedback).Error
}

// Delete 删除反馈
func (r *FeedbackRepository) Delete(id int) error {
	return r.DB.Delete(&models.UserFeedback{}, id).Error
}
