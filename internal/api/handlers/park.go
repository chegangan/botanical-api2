package handlers

import (
	"botanical-api2/internal/models"
	"botanical-api2/pkg/app"
	"botanical-api2/pkg/e"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPark 获取园区详情
// @Summary 获取园区详情
// @Description 根据ID获取园区详细信息
// @Tags 园区
// @Accept json
// @Produce json
// @Param id path int true "园区ID"
// @SUCCESS 200 {object} app.Response{data=models.Park} "获取成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "园区不存在"
// @Router /parks/{id} [get]
func (h *Handler) GetPark(c *gin.Context) {
	parkIDStr := c.Param("id")
	parkID, err := strconv.Atoi(parkIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的园区ID")
		return
	}

	park, err := h.ParkService.GetByID(parkID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取园区信息失败: "+err.Error())
		return
	}

	if park == nil {
		app.Error(c, e.NOT_FOUND, "园区不存在")
		return
	}

	app.SUCCESS(c, park)
}

// GetParks 获取园区列表
// @Summary 获取园区列表
// @Description 分页获取所有园区信息
// @Tags 园区
// @Accept json
// @Produce json
// @Param page query int false "页码(默认1)" default(1)
// @Param size query int false "每页条数(默认10)" default(10)
// @SUCCESS 200 {object} app.Response{data=app.PagedResult{list=[]models.Park}} "获取成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Router /parks [get]
func (h *Handler) GetParks(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

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

	parks, total, err := h.ParkService.GetAll(page, pageSize)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取园区列表失败: "+err.Error())
		return
	}

	app.SUCCESS(c, app.PagedResult{
		List:      parks,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: (total + pageSize - 1) / pageSize,
	})
}

// CreatePark 创建园区
// @Summary 创建园区
// @Description 创建新的园区信息
// @Tags 园区
// @Accept json
// @Produce json
// @Param park body models.Park true "园区信息"
// @SUCCESS 200 {object} app.Response{data=models.Park} "创建成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 500 {object} app.Response{data=string} "创建失败"
// @Security ApiKeyAuth
// @Router /parks [post]
func (h *Handler) CreatePark(c *gin.Context) {
	var park models.Park

	if err := c.ShouldBindJSON(&park); err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的请求参数: "+err.Error())
		return
	}

	if park.Src == "" {
		app.Error(c, e.BAD_REQUEST, "园区图片链接不能为空")
		return
	}

	if err := h.ParkService.Create(&park); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "创建园区失败: "+err.Error())
		return
	}

	app.SUCCESS(c, park)
}

// UpdatePark 更新园区信息
// @Summary 更新园区信息
// @Description 更新指定园区的信息
// @Tags 园区
// @Accept json
// @Produce json
// @Param id path int true "园区ID"
// @Param park body models.Park true "园区信息"
// @SUCCESS 200 {object} app.Response{data=models.Park} "更新成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "园区不存在"
// @Failure 500 {object} app.Response{data=string} "更新失败"
// @Security ApiKeyAuth
// @Router /parks/{id} [put]
func (h *Handler) UpdatePark(c *gin.Context) {
	parkIDStr := c.Param("id")
	parkID, err := strconv.Atoi(parkIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的园区ID")
		return
	}

	// 检查园区是否存在
	existingPark, err := h.ParkService.GetByID(parkID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取园区信息失败: "+err.Error())
		return
	}
	if existingPark == nil {
		app.Error(c, e.NOT_FOUND, "园区不存在")
		return
	}

	// 绑定请求体
	var park models.Park
	if err := c.ShouldBindJSON(&park); err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的请求参数: "+err.Error())
		return
	}

	// 确保ID正确
	park.ID = parkID

	if park.Src == "" {
		app.Error(c, e.BAD_REQUEST, "园区图片链接不能为空")
		return
	}

	if err := h.ParkService.Update(&park); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "更新园区失败: "+err.Error())
		return
	}

	app.SUCCESS(c, park)
}

// DeletePark 删除园区
// @Summary 删除园区
// @Description 删除指定园区
// @Tags 园区
// @Accept json
// @Produce json
// @Param id path int true "园区ID"
// @SUCCESS 200 {object} app.Response{data=map[string]interface{}} "删除成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "园区不存在"
// @Failure 500 {object} app.Response{data=string} "删除失败"
// @Security ApiKeyAuth
// @Router /parks/{id} [delete]
func (h *Handler) DeletePark(c *gin.Context) {
	parkIDStr := c.Param("id")
	parkID, err := strconv.Atoi(parkIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的园区ID")
		return
	}

	// 检查园区是否存在
	existingPark, err := h.ParkService.GetByID(parkID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取园区信息失败: "+err.Error())
		return
	}
	if existingPark == nil {
		app.Error(c, e.NOT_FOUND, "园区不存在")
		return
	}

	if err := h.ParkService.Delete(parkID); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "删除园区失败: "+err.Error())
		return
	}

	app.SUCCESS(c, gin.H{"message": "删除成功"})
}

// UploadParkImage 上传园区图片
// @Summary 上传园区图片
// @Description 为指定园区上传图片
// @Tags 园区,图片
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "园区ID"
// @Param image formData file true "图片文件"
// @SUCCESS 200 {object} app.Response{data=map[string]string} "上传成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "园区不存在"
// @Failure 500 {object} app.Response{data=string} "上传失败"
// @Security ApiKeyAuth
// @Router /parks/{id}/images [post]
func (h *Handler) UploadParkImage(c *gin.Context) {
	parkIDStr := c.Param("id")
	parkID, err := strconv.Atoi(parkIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的园区ID")
		return
	}

	// 获取上传文件
	file, err := c.FormFile("image")
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "请选择要上传的图片")
		return
	}

	// 上传图片并获取URL
	url, err := h.ParkService.UploadParkImage(parkID, file)
	if err != nil {
		app.Error(c, e.ERROR_UPLOAD_FAILED, "上传图片失败: "+err.Error())
		return
	}

	// 获取当前园区信息
	park, err := h.ParkService.GetByID(parkID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取园区信息失败: "+err.Error())
		return
	}

	// 更新园区图片链接
	park.Src = url

	// 更新园区信息
	if err := h.ParkService.Update(park); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "更新园区图片信息失败: "+err.Error())
		return
	}

	app.SUCCESS(c, gin.H{
		"url": url,
	})
}

// AddPlantToPark 将植物添加到园区
// @Summary 将植物添加到园区
// @Description 建立园区与植物的关联关系
// @Tags 园区,植物
// @Accept json
// @Produce json
// @Param id path int true "园区ID"
// @Param plant_id path int true "植物ID"
// @SUCCESS 200 {object} app.Response{data=map[string]interface{}} "添加成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "园区或植物不存在"
// @Failure 500 {object} app.Response{data=string} "添加失败"
// @Security ApiKeyAuth
// @Router /parks/{park_id}/plants/{plant_id} [post]
func (h *Handler) AddPlantToPark(c *gin.Context) {
	parkIDStr := c.Param("id")
	plantIDStr := c.Param("plant_id")

	parkID, err := strconv.Atoi(parkIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的园区ID")
		return
	}

	plantID, err := strconv.Atoi(plantIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的植物ID")
		return
	}

	err = h.ParkService.AddPlantToPark(parkID, plantID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "添加植物到园区失败: "+err.Error())
		return
	}

	app.SUCCESS(c, gin.H{"message": "添加成功"})
}

// RemovePlantFromPark 从园区移除植物
// @Summary 从园区移除植物
// @Description 解除园区与植物的关联关系
// @Tags 园区,植物
// @Accept json
// @Produce json
// @Param id path int true "园区ID"
// @Param plant_id path int true "植物ID"
// @SUCCESS 200 {object} app.Response{data=map[string]interface{}} "移除成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 500 {object} app.Response{data=string} "移除失败"
// @Security ApiKeyAuth
// @Router /parks/{park_id}/plants/{plant_id} [delete]
func (h *Handler) RemovePlantFromPark(c *gin.Context) {
	parkIDStr := c.Param("id")
	plantIDStr := c.Param("plant_id")

	parkID, err := strconv.Atoi(parkIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的园区ID")
		return
	}

	plantID, err := strconv.Atoi(plantIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的植物ID")
		return
	}

	err = h.ParkService.RemovePlantFromPark(parkID, plantID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "从园区移除植物失败: "+err.Error())
		return
	}

	app.SUCCESS(c, gin.H{"message": "移除成功"})
}
