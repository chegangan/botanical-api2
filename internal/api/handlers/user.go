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
// @Summary 获取用户信息
// @Description 根据用户ID获取用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} app.Response{data=models.User} "操作成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 10001 {object} app.Response{data=string} "用户不存在"
// @Security ApiKeyAuth
// @Router /users/{id} [get]
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
// @Summary 创建新用户
// @Description 创建一个新的用户账号并返回用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body models.User true "用户信息"
// @Success 200 {object} app.Response{data=map[string]interface{}} "操作成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 10003 {object} app.Response{data=string} "创建用户失败"
// @Security ApiKeyAuth
// @Router /users [post]
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
// @Summary 更新用户信息
// @Description 根据用户ID更新用户信息，支持部分字段更新
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param user body map[string]interface{} true "要更新的用户字段"
// @Success 200 {object} app.Response{data=map[string]interface{}} "操作成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 10001 {object} app.Response{data=string} "用户不存在"
// @Failure 10004 {object} app.Response{data=string} "更新用户失败"
// @Security ApiKeyAuth
// @Router /users/{id} [put]
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
// @Summary 删除用户
// @Description 根据用户ID删除指定用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} app.Response{data=map[string]interface{}} "操作成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 10005 {object} app.Response{data=string} "删除用户失败"
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
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
// @Summary 修改用户密码
// @Description 修改指定用户的登录密码，需要验证原密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body dto.PasswordChangeRequest true "密码修改信息"
// @Success 200 {object} app.Response{data=map[string]interface{}} "操作成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 10004 {object} app.Response{data=string} "更新用户失败"
// @Failure 10007 {object} app.Response{data=string} "未授权访问"
// @Security ApiKeyAuth
// @Router /users/{id}/password [put]
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
