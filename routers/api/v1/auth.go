package v1

import (
	"botanical-api2/models"
	"botanical-api2/pkg/app"
	"botanical-api2/pkg/e"

	"github.com/gin-gonic/gin"
)

// RegisterUser 处理用户注册
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		app.Error(c, e.ErrorInvalidParams, err.Error())
		return
	}

	if err := userService.RegisterUser(&user); err != nil {
		app.Error(c, e.ErrorUserCreate, err.Error())
		return
	}

	app.Success(c, map[string]interface{}{
		"message": "用户注册成功",
		"user":    user,
	})
}

// LoginUser 处理用户登录
func LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		app.Error(c, e.ErrorInvalidParams, err.Error())
		return
	}

	// 这里的passwordhash其实是用户输入的密码，是明文
	token, err := userService.LoginUser(user.Phone, user.PasswordHash)
	if err != nil {
		app.Error(c, e.ErrorUserPasswordIncorrect, err.Error())
		return
	}

	app.Success(c, map[string]interface{}{
		"token": token,
	})
}
