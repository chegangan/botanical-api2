package handlers

import (
	"botanical-api2/internal/models"
	"botanical-api2/pkg/app"
	"botanical-api2/pkg/e"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 以下是植物相关处理方法，集成到主Handler中

// GetPlant 获取植物详情
// @Summary 获取植物详情
// @Description 根据ID获取植物详细信息
// @Tags 植物
// @Accept json
// @Produce json
// @Param id path int true "植物ID"
// @SUCCESS 200 {object} app.Response{data=models.PlantVegetation} "获取成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "植物不存在"
// @Router /plants/{id} [get]
func (h *Handler) GetPlant(c *gin.Context) {
	plantIDStr := c.Param("id")
	plantID, err := strconv.Atoi(plantIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的植物ID")
		return
	}

	plant, err := h.PlantService.GetByID(plantID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取植物信息失败: "+err.Error())
		return
	}

	if plant == nil {
		app.Error(c, e.NOT_FOUND, "植物不存在")
		return
	}

	app.SUCCESS(c, plant)
}

// GetPlants 获取植物列表
// @Summary 获取植物列表
// @Description 分页获取所有植物信息
// @Tags 植物
// @Accept json
// @Produce json
// @Param page query int false "页码(默认1)" default(1)
// @Param size query int false "每页条数(默认10)" default(10)
// @Param class_id query int false "分类ID(可选)"
// @SUCCESS 200 {object} app.Response{data=app.PagedResult{list=[]models.PlantVegetation}} "获取成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Router /plants [get]
func (h *Handler) GetPlants(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	classIDStr := c.Query("class_id")

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

	var plants []models.PlantVegetation
	var total int

	if classIDStr != "" {
		classID, err := strconv.Atoi(classIDStr)
		if err != nil {
			app.Error(c, e.BAD_REQUEST, "无效的分类ID")
			return
		}
		var err2 error
		plants, total, err2 = h.PlantService.GetByClass(classID, page, pageSize)
		if err2 != nil {
			app.Error(c, e.INTERNAL_SERVER, "获取植物列表失败: "+err2.Error())
			return
		}
	} else {
		plants, total, err = h.PlantService.GetAll(page, pageSize)
		if err != nil {
			app.Error(c, e.INTERNAL_SERVER, "获取植物列表失败: "+err.Error())
			return
		}
	}

	app.SUCCESS(c, app.PagedResult{
		List:      plants,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: (total + pageSize - 1) / pageSize,
	})
}

// CreatePlant 创建植物
// @Summary 创建植物
// @Description 创建新的植物信息
// @Tags 植物
// @Accept json
// @Produce json
// @Param plant body models.PlantVegetation true "植物信息"
// @SUCCESS 200 {object} app.Response{data=models.PlantVegetation} "创建成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 500 {object} app.Response{data=string} "创建失败"
// @Security ApiKeyAuth
// @Router /plants [post]
func (h *Handler) CreatePlant(c *gin.Context) {
	var plant models.PlantVegetation

	if err := c.ShouldBindJSON(&plant); err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的请求参数: "+err.Error())
		return
	}

	if err := h.PlantService.Create(&plant); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "创建植物失败: "+err.Error())
		return
	}

	app.SUCCESS(c, plant)
}

// UpdatePlant 更新植物信息
// @Summary 更新植物信息
// @Description 更新指定植物的信息
// @Tags 植物
// @Accept json
// @Produce json
// @Param id path int true "植物ID"
// @Param plant body models.PlantVegetation true "植物信息"
// @SUCCESS 200 {object} app.Response{data=models.PlantVegetation} "更新成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "植物不存在"
// @Failure 500 {object} app.Response{data=string} "更新失败"
// @Security ApiKeyAuth
// @Router /plants/{id} [put]
func (h *Handler) UpdatePlant(c *gin.Context) {
	plantIDStr := c.Param("id")
	plantID, err := strconv.Atoi(plantIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的植物ID")
		return
	}

	// 检查植物是否存在
	existingPlant, err := h.PlantService.GetByID(plantID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取植物信息失败: "+err.Error())
		return
	}
	if existingPlant == nil {
		app.Error(c, e.NOT_FOUND, "植物不存在")
		return
	}

	// 绑定请求体
	var plant models.PlantVegetation
	if err := c.ShouldBindJSON(&plant); err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的请求参数: "+err.Error())
		return
	}

	// 确保ID正确
	plant.ID = plantID

	if err := h.PlantService.Update(&plant); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "更新植物失败: "+err.Error())
		return
	}

	app.SUCCESS(c, plant)
}

// DeletePlant 删除植物
// @Summary 删除植物
// @Description 删除指定植物
// @Tags 植物
// @Accept json
// @Produce json
// @Param id path int true "植物ID"
// @SUCCESS 200 {object} app.Response{data=map[string]interface{}} "删除成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "植物不存在"
// @Failure 500 {object} app.Response{data=string} "删除失败"
// @Security ApiKeyAuth
// @Router /plants/{id} [delete]
func (h *Handler) DeletePlant(c *gin.Context) {
	plantIDStr := c.Param("id")
	plantID, err := strconv.Atoi(plantIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的植物ID")
		return
	}

	// 检查植物是否存在
	existingPlant, err := h.PlantService.GetByID(plantID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取植物信息失败: "+err.Error())
		return
	}
	if existingPlant == nil {
		app.Error(c, e.NOT_FOUND, "植物不存在")
		return
	}

	if err := h.PlantService.Delete(plantID); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "删除植物失败: "+err.Error())
		return
	}

	app.SUCCESS(c, gin.H{"message": "删除成功"})
}

// GetParkPlants 获取园区植物列表
// @Summary 获取园区植物列表
// @Description 获取指定园区内的所有植物
// @Tags 植物,园区
// @Accept json
// @Produce json
// @Param id path int true "园区ID"
// @Param page query int false "页码(默认1)" default(1)
// @Param size query int false "每页条数(默认10)" default(10)
// @SUCCESS 200 {object} app.Response{data=app.PagedResult{list=[]models.PlantVegetation}} "获取成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Router /parks/{park_id}/plants [get]
func (h *Handler) GetParkPlants(c *gin.Context) {
	parkIDStr := c.Param("id")
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	parkID, err := strconv.Atoi(parkIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的园区ID")
		return
	}

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

	plants, total, err := h.PlantService.GetByParkID(parkID, page, pageSize)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取园区植物列表失败: "+err.Error())
		return
	}

	app.SUCCESS(c, app.PagedResult{
		List:      plants,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: (total + pageSize - 1) / pageSize,
	})
}

// UploadPlantImage 上传植物图片
// @Summary 上传植物图片
// @Description 为指定植物上传图片
// @Tags 植物,图片
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "植物ID"
// @Param type query string true "图片类型(main/list)" Enums(main, list)
// @Param image formData file true "图片文件"
// @SUCCESS 200 {object} app.Response{data=map[string]string} "上传成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "植物不存在"
// @Failure 500 {object} app.Response{data=string} "上传失败"
// @Security ApiKeyAuth
// @Router /plants/{id}/images [post]
func (h *Handler) UploadPlantImage(c *gin.Context) {
	plantIDStr := c.Param("id")
	imageType := c.Query("type")

	if imageType != "main" && imageType != "list" {
		app.Error(c, e.BAD_REQUEST, "图片类型必须为main或list")
		return
	}

	plantID, err := strconv.Atoi(plantIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的植物ID")
		return
	}

	// 获取上传文件
	file, err := c.FormFile("image")
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "请选择要上传的图片")
		return
	}

	// 上传图片并获取URL
	url, err := h.PlantService.UploadPlantImage(plantID, file)
	if err != nil {
		app.Error(c, e.ERROR_UPLOAD_FAILED, "上传图片失败: "+err.Error())
		return
	}

	// 获取当前植物信息
	plant, err := h.PlantService.GetByID(plantID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取植物信息失败: "+err.Error())
		return
	}

	// 根据图片类型更新相应字段
	if imageType == "main" {
		plant.Src = url
	} else if imageType == "list" {
		// 处理多图URL列表
		var srcList []string
		if plant.SrcList != "" {
			srcList = strings.Split(plant.SrcList, ",")
		}
		srcList = append(srcList, url)
		plant.SrcList = strings.Join(srcList, ",")
	}

	// 更新植物信息
	if err := h.PlantService.Update(plant); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "更新植物图片信息失败: "+err.Error())
		return
	}

	app.SUCCESS(c, gin.H{
		"url":        url,
		"image_type": imageType,
	})
}
