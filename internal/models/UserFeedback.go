package models

import "time"

// UserFeedback 用户反馈模型
// @Description 用户提交的反馈信息
type UserFeedback struct {
	ID              int       `json:"id" gorm:"primary_key;auto_increment;unique"`                                               // 反馈ID
	FeedbackType    string    `json:"feedback_type" gorm:"column:feedback_type;size:255"`                                        // 反馈类型
	UserID          int       `json:"user_id" gorm:"column:user_id;index"`                                                       // 用户ID，外键
	FeedbackContent string    `json:"feedback_content" gorm:"column:feedback_content;type:text"`                                 // 反馈内容
	FeedbackSrcList string    `json:"feedback_src_list" gorm:"column:feedback_src_list;type:text"`                               // 图片列表，逗号分隔的URL
	FeedbackStatus  string    `json:"feedback_status" gorm:"column:feedback_status;size:50;default:'pending'"`                   // 反馈状态：pending(待处理), processing(处理中), resolved(已解决)
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`                             // 创建时间
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新时间
}

// TableName 设置表名
func (UserFeedback) TableName() string {
	return "user_feedback"
}
