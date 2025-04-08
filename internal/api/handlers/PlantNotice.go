package handlers

import (
	"botanical-api2/internal/models"
	"botanical-api2/pkg/app"
	"botanical-api2/pkg/e"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPlantNotice 获取植物资讯详情
// @Summary 获取植物资讯详情
// @Description 根据ID获取植物资讯详细信息
// @Tags 植物资讯
// @Accept json
// @Produce json
// @Param id path int true "资讯ID"
// @SUCCESS 200 {object} app.Response{data=models.PlantNotice} "获取成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "资讯不存在"
// @Router /notices/{id} [get]
func (h *Handler) GetPlantNotice(c *gin.Context) {
	noticeIDStr := c.Param("id")
	noticeID, err := strconv.Atoi(noticeIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的资讯ID")
		return
	}

	// 获取资讯信息
	notice, err := h.PlantNoticeService.GetByID(noticeID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取资讯信息失败: "+err.Error())
		return
	}

	if notice == nil {
		app.Error(c, e.NOT_FOUND, "资讯不存在")
		return
	}

	// 增加浏览量
	go h.PlantNoticeService.IncrementSkimSort(noticeID)

	app.SUCCESS(c, notice)
}

// GetPlantNotices 获取植物资讯列表
// @Summary 获取植物资讯列表
// @Description 分页获取所有植物资讯信息
// @Tags 植物资讯
// @Accept json
// @Produce json
// @Param page query int false "页码(默认1)" default(1)
// @Param size query int false "每页条数(默认10)" default(10)
// @Param recommend query int false "是否推荐(0不推荐/1推荐，不传则获取全部)"
// @SUCCESS 200 {object} app.Response{data=app.PagedResult{list=[]models.PlantNotice}} "获取成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Router /notices [get]
func (h *Handler) GetPlantNotices(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	recommendStr := c.Query("recommend")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		app.Error(c, e.BAD_REQUEST, "无效的页码")
		return
	}

	pageSize, err := strconv.Atoi(sizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		app.Error(c, e.BAD_REQUEST, "无效的每页条数")
		return
	}

	var isRecommend *int
	if recommendStr != "" {
		recommendVal, err := strconv.Atoi(recommendStr)
		if err != nil || (recommendVal != 0 && recommendVal != 1) {
			app.Error(c, e.BAD_REQUEST, "推荐参数应为0或1")
			return
		}
		isRecommend = &recommendVal
	}

	notices, total, err := h.PlantNoticeService.GetAll(page, pageSize, isRecommend)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取资讯列表失败: "+err.Error())
		return
	}

	app.SUCCESS(c, app.PagedResult{
		List:      notices,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: (total + pageSize - 1) / pageSize,
	})
}

// CreatePlantNotice 创建植物资讯
// @Summary 创建植物资讯
// @Description 创建新的植物资讯信息
// @Tags 植物资讯
// @Accept json
// @Produce json
// @Param notice body models.PlantNotice true "资讯信息"
// @SUCCESS 200 {object} app.Response{data=models.PlantNotice} "创建成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 500 {object} app.Response{data=string} "创建失败"
// @Security ApiKeyAuth
// @Router /notices [post]
func (h *Handler) CreatePlantNotice(c *gin.Context) {
	var notice models.PlantNotice

	if err := c.ShouldBindJSON(&notice); err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的请求参数: "+err.Error())
		return
	}

	if notice.Title == "" {
		app.Error(c, e.BAD_REQUEST, "资讯标题不能为空")
		return
	}

	// 设置默认值
	if notice.IsRecommend != 0 && notice.IsRecommend != 1 {
		notice.IsRecommend = 0 // 默认不推荐
	}

	if err := h.PlantNoticeService.Create(&notice); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "创建资讯失败: "+err.Error())
		return
	}

	app.SUCCESS(c, notice)
}

// UpdatePlantNotice 更新植物资讯信息
// @Summary 更新植物资讯信息
// @Description 更新指定植物资讯的信息
// @Tags 植物资讯
// @Accept json
// @Produce json
// @Param id path int true "资讯ID"
// @Param notice body models.PlantNotice true "资讯信息"
// @SUCCESS 200 {object} app.Response{data=models.PlantNotice} "更新成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "资讯不存在"
// @Failure 500 {object} app.Response{data=string} "更新失败"
// @Security ApiKeyAuth
// @Router /notices/{id} [put]
func (h *Handler) UpdatePlantNotice(c *gin.Context) {
	noticeIDStr := c.Param("id")
	noticeID, err := strconv.Atoi(noticeIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的资讯ID")
		return
	}

	// 检查资讯是否存在
	existingNotice, err := h.PlantNoticeService.GetByID(noticeID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取资讯信息失败: "+err.Error())
		return
	}
	if existingNotice == nil {
		app.Error(c, e.NOT_FOUND, "资讯不存在")
		return
	}

	// 绑定请求体
	var notice models.PlantNotice
	if err := c.ShouldBindJSON(&notice); err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的请求参数: "+err.Error())
		return
	}

	// 确保ID正确
	notice.ID = noticeID

	if notice.Title == "" {
		app.Error(c, e.BAD_REQUEST, "资讯标题不能为空")
		return
	}

	if notice.IsRecommend != 0 && notice.IsRecommend != 1 {
		notice.IsRecommend = existingNotice.IsRecommend // 保持原值
	}

	if err := h.PlantNoticeService.Update(&notice); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "更新资讯失败: "+err.Error())
		return
	}

	app.SUCCESS(c, notice)
}

// DeletePlantNotice 删除植物资讯
// @Summary 删除植物资讯
// @Description 删除指定植物资讯
// @Tags 植物资讯
// @Accept json
// @Produce json
// @Param id path int true "资讯ID"
// @SUCCESS 200 {object} app.Response{data=map[string]interface{}} "删除成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "资讯不存在"
// @Failure 500 {object} app.Response{data=string} "删除失败"
// @Security ApiKeyAuth
// @Router /notices/{id} [delete]
func (h *Handler) DeletePlantNotice(c *gin.Context) {
	noticeIDStr := c.Param("id")
	noticeID, err := strconv.Atoi(noticeIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的资讯ID")
		return
	}

	// 检查资讯是否存在
	existingNotice, err := h.PlantNoticeService.GetByID(noticeID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取资讯信息失败: "+err.Error())
		return
	}
	if existingNotice == nil {
		app.Error(c, e.NOT_FOUND, "资讯不存在")
		return
	}

	if err := h.PlantNoticeService.Delete(noticeID); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "删除资讯失败: "+err.Error())
		return
	}

	app.SUCCESS(c, gin.H{"message": "删除成功"})
}

// TogglePlantNoticeRecommend 切换资讯推荐状态
// @Summary 切换资讯推荐状态
// @Description 切换指定资讯的推荐状态(0不推荐/1推荐)
// @Tags 植物资讯
// @Accept json
// @Produce json
// @Param id path int true "资讯ID"
// @Param recommend body map[string]int true "推荐状态" example={"is_recommend":1}
// @SUCCESS 200 {object} app.Response{data=map[string]interface{}} "设置成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "资讯不存在"
// @Failure 500 {object} app.Response{data=string} "设置失败"
// @Security ApiKeyAuth
// @Router /notices/{id}/recommend [put]
func (h *Handler) TogglePlantNoticeRecommend(c *gin.Context) {
	noticeIDStr := c.Param("id")
	noticeID, err := strconv.Atoi(noticeIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的资讯ID")
		return
	}

	// 检查资讯是否存在
	existingNotice, err := h.PlantNoticeService.GetByID(noticeID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取资讯信息失败: "+err.Error())
		return
	}
	if existingNotice == nil {
		app.Error(c, e.NOT_FOUND, "资讯不存在")
		return
	}

	var reqBody struct {
		IsRecommend int `json:"is_recommend"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的请求参数: "+err.Error())
		return
	}

	if reqBody.IsRecommend != 0 && reqBody.IsRecommend != 1 {
		app.Error(c, e.BAD_REQUEST, "推荐状态必须为0或1")
		return
	}

	if err := h.PlantNoticeService.ToggleRecommend(noticeID, reqBody.IsRecommend); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "设置推荐状态失败: "+err.Error())
		return
	}

	app.SUCCESS(c, gin.H{
		"message":      "设置成功",
		"is_recommend": reqBody.IsRecommend,
	})
}

// UploadPlantNoticeImage 上传资讯图片
// @Summary 上传资讯图片
// @Description 为指定资讯上传图片
// @Tags 植物资讯,图片
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "资讯ID"
// @Param image formData file true "图片文件"
// @SUCCESS 200 {object} app.Response{data=map[string]string} "上传成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "资讯不存在"
// @Failure 500 {object} app.Response{data=string} "上传失败"
// @Security ApiKeyAuth
// @Router /notices/{id}/images [post]
func (h *Handler) UploadPlantNoticeImage(c *gin.Context) {
	noticeIDStr := c.Param("id")
	noticeID, err := strconv.Atoi(noticeIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的资讯ID")
		return
	}

	// 获取上传文件
	file, err := c.FormFile("image")
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "请选择要上传的图片")
		return
	}

	// 上传图片并获取URL
	url, err := h.PlantNoticeService.UploadNoticeImage(noticeID, file)
	if err != nil {
		app.Error(c, e.ERROR_UPLOAD_FAILED, "上传图片失败: "+err.Error())
		return
	}

	// 获取当前资讯信息
	notice, err := h.PlantNoticeService.GetByID(noticeID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取资讯信息失败: "+err.Error())
		return
	}

	// 更新资讯图片链接
	notice.Src = url

	// 更新资讯信息
	if err := h.PlantNoticeService.Update(notice); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "更新资讯图片信息失败: "+err.Error())
		return
	}

	app.SUCCESS(c, gin.H{
		"url": url,
	})
}
