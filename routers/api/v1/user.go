package v1

import (
	"botanical-api2/models"
	"botanical-api2/pkg/app"
	"botanical-api2/pkg/e"
	"botanical-api2/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserService 实例，负责处理用户相关业务逻辑
var userService *service.UserService

// SetUserService 设置用户服务实例
func SetUserService(service *service.UserService) {
	userService = service
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		app.Error(c, e.ErrorInvalidParams, err.Error())
		return
	}

	if err := userService.CreateUser(&user); err != nil {
		app.Error(c, e.ErrorUserCreate, err.Error())
		return
	}

	app.Success(c, map[string]interface{}{
		"message": "用户创建成功",
		"user":    user,
	})
}

// GetUser 获取用户
func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	// 字符串转整数
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.Error(c, e.ErrorInvalidParams, "无效的用户ID")
		return
	}

	user, err := userService.GetUserByID(id)
	if err != nil {
		app.Error(c, e.ErrorUserNotFound, "用户不存在")
		return
	}
	app.Success(c, user)
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	// 字符串转整数
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.Error(c, e.ErrorInvalidParams, "无效的用户ID")
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		app.Error(c, e.ErrorInvalidParams, err.Error())
		return
	}
	user.ID = id
	if err := userService.UpdateUser(&user); err != nil {
		app.Error(c, e.ErrorUserUpdate, err.Error())
		return
	}
	app.Success(c, map[string]interface{}{
		"message": "用户更新成功",
		"user":    user,
	})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	// 字符串转整数
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.Error(c, e.ErrorInvalidParams, "无效的用户ID")
		return
	}

	if err := userService.DeleteUser(id); err != nil {
		app.Error(c, e.ErrorUserDelete, err.Error())
		return
	}
	app.Success(c, map[string]interface{}{
		"message": "用户删除成功",
	})
}
