package dto

// CreateFeedbackRequest 创建反馈请求参数
type CreateFeedbackRequest struct {
	Content string `json:"content" binding:"required" example:"反馈内容"`
	Type    string `json:"type" binding:"required" example:"功能建议"`
}

// UpdateFeedbackStatusRequest 更新反馈状态请求参数
type UpdateFeedbackStatusRequest struct {
	Status string `json:"status" binding:"required" example:"processing"`
}
