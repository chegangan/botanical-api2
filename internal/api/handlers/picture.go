package handlers

import (
	"botanical-api2/internal/models"
	"botanical-api2/pkg/app"
	"botanical-api2/pkg/e"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 以下是图片相关处理方法，集成到主Handler中

// UploadAvatar 上传头像
// @Summary 上传用户头像
// @Description 上传并更新当前登录用户的头像
// @Tags 用户头像
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "头像图片文件"
// @SUCCESS 200 {object} app.Response{data=models.UserAvatar} "上传成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 500 {object} app.Response{data=string} "上传失败"
// @Security ApiKeyAuth
// @Router /avatar [post]
func (h *Handler) UploadAvatar(c *gin.Context) {
	// 从上下文中获取用户（通过中间件设置）
	userObj, exists := c.Get("user")
	if !exists {
		app.Error(c, e.ERROR_PICTURE_NO_PERMISSION, "未授权")
		return
	}

	// 类型断言获取用户对象
	user, ok := userObj.(*models.User)
	if !ok {
		app.Error(c, e.INTERNAL_SERVER, "用户信息类型错误")
		return
	}

	// 使用用户ID
	userID := user.ID

	// 获取上传文件
	file, err := c.FormFile("avatar")
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "请选择要上传的头像")
		return
	}

	// 上传头像
	avatar, err := h.PictureService.UploadAvatar(userID, file)
	if err != nil {
		app.Error(c, e.ERROR_UPLOAD_FAILED, "上传头像失败: "+err.Error())
		return
	}

	app.SUCCESS(c, avatar)
}

// GetUserAvatar 获取用户头像
// @Summary 获取用户头像
// @Description 获取指定用户的头像信息
// @Tags 用户头像
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @SUCCESS 200 {object} app.Response{data=models.UserAvatar} "获取成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "头像不存在"
// @Router /users/{id}/avatar [get]
func (h *Handler) GetUserAvatar(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的用户ID")
		return
	}

	avatar, err := h.PictureService.GetUserAvatar(userID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取头像失败")
		return
	}

	if avatar == nil {
		app.Error(c, e.NOT_FOUND, "用户未设置头像")
		return
	}

	app.SUCCESS(c, avatar)
}

// UploadUserPicture 上传用户图片
// @Summary 上传用户图片
// @Description 上传当前登录用户的图片
// @Tags 用户图片
// @Accept multipart/form-data
// @Produce json
// @Param picture formData file true "图片文件"
// @SUCCESS 200 {object} app.Response{data=models.UserPicture} "上传成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 500 {object} app.Response{data=string} "上传失败"
// @Security ApiKeyAuth
// @Router /pictures [post]
func (h *Handler) UploadUserPicture(c *gin.Context) {
	// 从上下文中获取用户对象（通过中间件设置）
	userObj, exists := c.Get("user")
	if !exists {
		app.Error(c, e.ERROR_PICTURE_NO_PERMISSION, "未授权")
		return
	}

	// 类型断言获取用户对象
	user, ok := userObj.(*models.User)
	if !ok {
		app.Error(c, e.INTERNAL_SERVER, "用户信息类型错误")
		return
	}

	file, err := c.FormFile("picture")
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "请选择要上传的图片")
		return
	}

	picture, err := h.PictureService.UploadUserPicture(user.ID, file)
	if err != nil {
		app.Error(c, e.ERROR_UPLOAD_FAILED, "上传图片失败: "+err.Error())
		return
	}

	app.SUCCESS(c, picture)
}

// GetUserPictures 获取用户所有图片
// @Summary 获取用户所有图片
// @Description 获取指定用户的所有图片列表
// @Tags 用户图片
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @SUCCESS 200 {object} app.Response{data=[]models.UserPicture} "获取成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Router /users/{id}/pictures [get]
func (h *Handler) GetUserPictures(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的用户ID")
		return
	}

	pictures, err := h.PictureService.GetUserPictures(userID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取图片失败")
		return
	}

	app.SUCCESS(c, pictures)
}

// GetPicture 获取单张图片
// @Summary 获取单张图片
// @Description 根据图片ID获取详细信息
// @Tags 用户图片
// @Accept json
// @Produce json
// @Param id path int true "图片ID"
// @SUCCESS 200 {object} app.Response{data=models.UserPicture} "获取成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "图片不存在"
// @Router /pictures/{id} [get]
func (h *Handler) GetPicture(c *gin.Context) {
	pictureIDStr := c.Param("id")
	pictureID, err := strconv.Atoi(pictureIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的图片ID")
		return
	}

	picture, err := h.PictureService.GetPictureByID(pictureID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取图片失败")
		return
	}

	if picture == nil {
		app.Error(c, e.NOT_FOUND, "图片不存在")
		return
	}

	app.SUCCESS(c, picture)
}

// DeletePicture 删除图片
// @Summary 删除图片
// @Description 删除当前用户的图片
// @Tags 用户图片
// @Accept json
// @Produce json
// @Param id path int true "图片ID"
// @SUCCESS 200 {object} app.Response{data=map[string]interface{}} "删除成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 403 {object} app.Response{data=string} "无权删除该图片"
// @Failure 500 {object} app.Response{data=string} "删除失败"
// @Security ApiKeyAuth
// @Router /pictures/{id} [delete]
func (h *Handler) DeletePicture(c *gin.Context) {
	// 从上下文中获取用户对象（通过中间件设置）
	userObj, exists := c.Get("user")
	if !exists {
		app.Error(c, e.ERROR_PICTURE_NO_PERMISSION, "未授权")
		return
	}

	// 类型断言获取用户对象
	user, ok := userObj.(*models.User)
	if !ok {
		app.Error(c, e.INTERNAL_SERVER, "用户信息类型错误")
		return
	}

	pictureIDStr := c.Param("id")
	pictureID, err := strconv.Atoi(pictureIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的图片ID")
		return
	}

	err = h.PictureService.DeletePicture(user.ID, pictureID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "删除图片失败: "+err.Error())
		return
	}

	app.SUCCESS(c, gin.H{"message": "删除成功"})
}
