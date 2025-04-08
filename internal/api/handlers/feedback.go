package handlers

import (
	"botanical-api2/internal/dto"
	"botanical-api2/internal/models"
	"botanical-api2/pkg/app"
	"botanical-api2/pkg/e"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 以下是反馈相关处理方法，集成到主Handler中

// CreateFeedback 创建用户反馈
// @Summary 创建用户反馈
// @Description 提交用户反馈信息
// @Tags 用户反馈
// @Accept json
// @Produce json
// @Param data body dto.CreateFeedbackRequest true "反馈信息"
// @SUCCESS 200 {object} app.Response{data=models.UserFeedback} "提交成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 500 {object} app.Response{data=string} "提交失败"
// @Security ApiKeyAuth
// @Router /feedback [post]
func (h *Handler) CreateFeedback(c *gin.Context) {
	// 从上下文中获取用户
	userObj, exists := c.Get("user")
	if !exists {
		app.Error(c, e.UNAUTHORIZED, "未授权")
		return
	}

	// 类型断言获取用户对象
	user, ok := userObj.(*models.User)
	if !ok {
		app.Error(c, e.INTERNAL_SERVER, "用户信息类型错误")
		return
	}

	var req dto.CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的请求参数")
		return
	}

	// 创建反馈
	feedback, err := h.FeedbackService.CreateFeedback(user.ID, req.Content, req.Type)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "创建反馈失败: "+err.Error())
		return
	}

	app.SUCCESS(c, feedback)
}

// GetUserFeedbacks 获取用户所有反馈
// @Summary 获取用户所有反馈
// @Description 获取指定用户的所有反馈列表
// @Tags 用户反馈
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @SUCCESS 200 {object} app.Response{data=[]models.UserFeedback} "获取成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Router /users/{id}/feedbacks [get]
func (h *Handler) GetUserFeedbacks(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的用户ID")
		return
	}

	feedbacks, err := h.FeedbackService.GetUserFeedbacks(userID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取反馈失败")
		return
	}

	app.SUCCESS(c, feedbacks)
}

// GetFeedback 获取单个反馈
// @Summary 获取单个反馈
// @Description 根据反馈ID获取详细信息
// @Tags 用户反馈
// @Accept json
// @Produce json
// @Param id path int true "反馈ID"
// @SUCCESS 200 {object} app.Response{data=models.UserFeedback} "获取成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "反馈不存在"
// @Router /feedback/{id} [get]
func (h *Handler) GetFeedback(c *gin.Context) {
	feedbackIDStr := c.Param("id")
	feedbackID, err := strconv.Atoi(feedbackIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的反馈ID")
		return
	}

	feedback, err := h.FeedbackService.GetFeedbackByID(feedbackID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取反馈失败")
		return
	}

	if feedback == nil {
		app.Error(c, e.NOT_FOUND, "反馈不存在")
		return
	}

	app.SUCCESS(c, feedback)
}

// UpdateFeedbackStatus 更新反馈状态
// @Summary 更新反馈状态
// @Description 更新反馈的处理状态（管理员使用）
// @Tags 用户反馈
// @Accept json
// @Produce json
// @Param id path int true "反馈ID"
// @Param data body dto.UpdateFeedbackStatusRequest true "状态信息"
// @SUCCESS 200 {object} app.Response{data=models.UserFeedback} "更新成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 403 {object} app.Response{data=string} "无权限操作"
// @Failure 404 {object} app.Response{data=string} "反馈不存在"
// @Security ApiKeyAuth
// @Router /feedback/{id}/status [put]
func (h *Handler) UpdateFeedbackStatus(c *gin.Context) {
	// TODO: 这里应该添加管理员权限检查

	feedbackIDStr := c.Param("id")
	feedbackID, err := strconv.Atoi(feedbackIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的反馈ID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的请求参数")
		return
	}

	feedback, err := h.FeedbackService.UpdateFeedbackStatus(feedbackID, req.Status)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "更新反馈状态失败: "+err.Error())
		return
	}

	app.SUCCESS(c, feedback)
}

// DeleteFeedback 删除反馈
// @Summary 删除反馈
// @Description 删除用户反馈
// @Tags 用户反馈
// @Accept json
// @Produce json
// @Param id path int true "反馈ID"
// @SUCCESS 200 {object} app.Response{data=map[string]interface{}} "删除成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 403 {object} app.Response{data=string} "无权删除该反馈"
// @Failure 500 {object} app.Response{data=string} "删除失败"
// @Security ApiKeyAuth
// @Router /feedback/{id} [delete]
func (h *Handler) DeleteFeedback(c *gin.Context) {
	// 从上下文中获取用户对象
	userObj, exists := c.Get("user")
	if !exists {
		app.Error(c, e.UNAUTHORIZED, "未授权")
		return
	}

	// 类型断言获取用户对象
	user, ok := userObj.(*models.User)
	if !ok {
		app.Error(c, e.INTERNAL_SERVER, "用户信息类型错误")
		return
	}

	feedbackIDStr := c.Param("id")
	feedbackID, err := strconv.Atoi(feedbackIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的反馈ID")
		return
	}

	// TODO: 检查用户是否为管理员，此处简单实现
	isAdmin := false // 这应该从用户角色中获取

	err = h.FeedbackService.DeleteFeedback(user.ID, feedbackID, isAdmin)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "删除反馈失败: "+err.Error())
		return
	}

	app.SUCCESS(c, gin.H{"message": "删除成功"})
}

// GetAllFeedbacks 获取所有反馈（管理员使用）
// @Summary 获取所有反馈
// @Description 获取系统中所有用户的反馈（管理员使用）
// @Tags 用户反馈
// @Accept json
// @Produce json
// @SUCCESS 200 {object} app.Response{data=[]models.UserFeedback} "获取成功"
// @Failure 403 {object} app.Response{data=string} "无权限操作"
// @Security ApiKeyAuth
// @Router /admin/feedbacks [get]
func (h *Handler) GetAllFeedbacks(c *gin.Context) {
	// TODO: 这里应该添加管理员权限检查

	feedbacks, err := h.FeedbackService.GetAllFeedbacks()
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取反馈失败")
		return
	}

	app.SUCCESS(c, feedbacks)
}
