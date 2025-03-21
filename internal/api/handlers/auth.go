package handlers

import (
	"botanical-api2/internal/dto"
	"botanical-api2/internal/models"
	"botanical-api2/pkg/app"
	"botanical-api2/pkg/e"

	"github.com/gin-gonic/gin"
)

// RegisterUser 处理用户注册
// @Summary 注册新用户
// @Description 创建新用户账号并返回用户信息和访问令牌
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "用户注册信息"
// @Success 200 {object} app.Response{data=map[string]interface{}} "注册成功，返回用户信息和令牌"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 10002 {object} app.Response{data=string} "用户已存在"
// @Failure 10003 {object} app.Response{data=string} "创建用户失败"
// @Router /auth/register [post]
func (h *Handler) RegisterUser(c *gin.Context) {
	// 使用DTO接收请求
	var registerReq dto.RegisterRequest
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		app.Error(c, e.ErrorInvalidParams, err.Error())
		return
	}

	// 转换为模型对象
	user := &models.User{
		Username:     registerReq.Username,
		Phone:        registerReq.Phone,
		PasswordHash: registerReq.Password, // service层会处理加密
	}

	// 注册用户
	if err := h.UserService.RegisterUser(user); err != nil {
		app.Error(c, e.ErrorUserCreate, err.Error())
		return
	}

	// 注册成功后，立即生成令牌（与登录流程相同）
	token, userInfo, err := h.UserService.LoginUser(registerReq.Phone, registerReq.Password)
	if err != nil {
		// 虽然注册成功，但令牌生成失败，返回用户信息但不返回令牌
		app.Success(c, map[string]interface{}{
			"message": "用户注册成功，但无法生成令牌",
			"user": dto.UserSummary{
				ID:       user.ID,
				Username: user.Username,
				Phone:    user.Phone,
				UserRole: user.UserRole,
			},
		})
		return
	}

	// 使用统一的认证响应DTO
	response := dto.AuthResponse{
		Token: token,
		User: dto.UserSummary{
			ID:       userInfo.ID,
			Username: userInfo.Username,
			Phone:    userInfo.Phone,
			UserRole: userInfo.UserRole,
		},
	}

	app.Success(c, map[string]interface{}{
		"message": "用户注册成功",
		"token":   token,
		"user":    response.User,
	})
}

// LoginUser 处理用户登录
// @Summary 用户登录
// @Description 验证用户凭据并返回访问令牌
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录凭据"
// @Success 200 {object} app.Response{data=dto.AuthResponse} "登录成功，返回用户信息和令牌"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 10006 {object} app.Response{data=string} "用户名或密码不正确"
// @Router /auth/login [post]
func (h *Handler) LoginUser(c *gin.Context) {
	// 使用DTO接收请求
	var loginReq dto.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		app.Error(c, e.ErrorInvalidParams, err.Error())
		return
	}

	// 将请求中的password传给service层
	token, userInfo, err := h.UserService.LoginUser(loginReq.Phone, loginReq.Password)
	if err != nil {
		app.Error(c, e.ErrorUserPasswordIncorrect, err.Error())
		return
	}

	// 使用DTO构建响应
	response := dto.AuthResponse{
		Token: token,
		User: dto.UserSummary{
			ID:       userInfo.ID,
			Username: userInfo.Username,
			Phone:    userInfo.Phone,
			UserRole: userInfo.UserRole,
		},
	}

	app.Success(c, response)
}
