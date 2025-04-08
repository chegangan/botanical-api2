package service

import (
	"errors"
	"fmt"

	"botanical-api2/internal/models"
	"botanical-api2/internal/repository"
)

// FeedbackService 反馈服务接口
type FeedbackService struct {
	Repositories *repository.Repositories
}

// NewFeedbackService 创建反馈服务
func NewFeedbackService(repositories *repository.Repositories) *FeedbackService {
	return &FeedbackService{
		Repositories: repositories,
	}
}

// CreateFeedback 创建用户反馈
func (s *FeedbackService) CreateFeedback(userID int, content string, feedbackType string) (*models.UserFeedback, error) {
	if content == "" {
		return nil, errors.New("反馈内容不能为空")
	}

	feedback := &models.UserFeedback{
		UserID:          userID,
		FeedbackContent: content,
		FeedbackType:    feedbackType,
		FeedbackStatus:  "pending", // 默认状态为待处理
	}

	if err := s.Repositories.Feedback.Create(feedback); err != nil {
		return nil, fmt.Errorf("创建反馈失败: %w", err)
	}

	return feedback, nil
}

// GetFeedbackByID 获取单条反馈
func (s *FeedbackService) GetFeedbackByID(id int) (*models.UserFeedback, error) {
	return s.Repositories.Feedback.GetByID(id)
}

// GetUserFeedbacks 获取用户所有反馈
func (s *FeedbackService) GetUserFeedbacks(userID int) ([]models.UserFeedback, error) {
	return s.Repositories.Feedback.GetByUserID(userID)
}

// GetAllFeedbacks 获取所有反馈（管理员使用）
func (s *FeedbackService) GetAllFeedbacks() ([]models.UserFeedback, error) {
	return s.Repositories.Feedback.GetAll()
}

// UpdateFeedbackStatus 更新反馈状态
func (s *FeedbackService) UpdateFeedbackStatus(id int, status string) (*models.UserFeedback, error) {
	feedback, err := s.Repositories.Feedback.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取反馈信息失败: %w", err)
	}

	if feedback == nil {
		return nil, errors.New("反馈不存在")
	}

	// 更新状态
	feedback.FeedbackStatus = status
	if err := s.Repositories.Feedback.Update(feedback); err != nil {
		return nil, fmt.Errorf("更新反馈状态失败: %w", err)
	}

	return feedback, nil
}

// DeleteFeedback 删除反馈
func (s *FeedbackService) DeleteFeedback(userID int, feedbackID int, isAdmin bool) error {
	// 获取反馈信息
	feedback, err := s.Repositories.Feedback.GetByID(feedbackID)
	if err != nil {
		return fmt.Errorf("获取反馈信息失败: %w", err)
	}

	if feedback == nil {
		return errors.New("反馈不存在")
	}

	// 检查权限：只有反馈创建者或管理员可以删除
	if feedback.UserID != userID && !isAdmin {
		return errors.New("无权删除该反馈")
	}

	// 删除数据库记录
	return s.Repositories.Feedback.Delete(feedbackID)
}
