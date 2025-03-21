package handlers

import (
	"botanical-api2/internal/dto"
	"botanical-api2/internal/models"
	"botanical-api2/pkg/app"
	"botanical-api2/pkg/e"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUser 获取用户
func (h *Handler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	// 字符串转整数
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.Error(c, e.ErrorInvalidParams, "无效的用户ID")
		return
	}

	user, err := h.UserService.GetUserByID(id)
	if err != nil {
		app.Error(c, e.ErrorUserNotFound, "用户不存在")
		return
	}
	app.Success(c, user)
}

// CreateUser 创建用户
func (h *Handler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		app.Error(c, e.ErrorInvalidParams, err.Error())
		return
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		app.Error(c, e.ErrorUserCreate, err.Error())
		return
	}

	app.Success(c, map[string]interface{}{
		"message": "用户创建成功",
		"user":    user,
	})
}

// UpdateUser 更新用户
func (h *Handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	// 字符串转整数
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.Error(c, e.ErrorInvalidParams, "无效的用户ID")
		return
	}

	// 使用map来接收动态字段
	var updateFields map[string]interface{}
	if err := c.ShouldBindJSON(&updateFields); err != nil {
		app.Error(c, e.ErrorInvalidParams, err.Error())
		return
	}

	// 创建用户对象，只设置ID
	user := &models.User{ID: id}

	// 调用服务层进行动态更新
	if err := h.UserService.UpdateUser(user, updateFields); err != nil {
		app.Error(c, e.ErrorUserUpdate, err.Error())
		return
	}

	// 获取更新后的完整用户信息
	updatedUser, _ := h.UserService.GetUserByID(id)

	app.Success(c, map[string]interface{}{
		"message": "用户更新成功",
		"user":    updatedUser,
	})
}

// DeleteUser 删除用户
func (h *Handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	// 字符串转整数
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.Error(c, e.ErrorInvalidParams, "无效的用户ID")
		return
	}

	if err := h.UserService.DeleteUser(id); err != nil {
		app.Error(c, e.ErrorUserDelete, err.Error())
		return
	}
	app.Success(c, map[string]interface{}{
		"message": "用户删除成功",
	})
}

// ChangePassword 修改用户密码
func (h *Handler) ChangePassword(c *gin.Context) {
	// 获取用户ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.Error(c, e.ErrorInvalidParams, "无效的用户ID")
		return
	}

	// 从JWT中获取当前用户
	currentUser, exists := c.Get("user")
	if !exists {
		app.Error(c, e.ErrorUserUnauthorized, "用户未认证")
		return
	}

	// 类型断言
	user, ok := currentUser.(*models.User)
	if !ok {
		app.Error(c, e.ErrorUserUnauthorized, "无效的用户信息")
		return
	}

	// 检查是否为用户本人或管理员
	if user.ID != id && user.UserRole != 9 {
		app.Error(c, e.ErrorUserUnauthorized, "无权修改其他用户密码")
		return
	}

	// 绑定请求数据到DTO
	var req dto.PasswordChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, e.ErrorInvalidParams, err.Error())
		return
	}

	// 调用服务层修改密码
	if err := h.UserService.ChangePassword(id, req.OldPassword, req.NewPassword); err != nil {
		app.Error(c, e.ErrorUserUpdate, err.Error())
		return
	}

	app.Success(c, map[string]interface{}{
		"message": "密码修改成功",
	})
}
