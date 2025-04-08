package handlers

import (
	"botanical-api2/internal/models"
	"botanical-api2/pkg/app"
	"botanical-api2/pkg/e"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetPlantClass 获取植物分类详情
// @Summary 获取植物分类详情
// @Description 根据ID获取植物分类详细信息
// @Tags 植物分类
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @SUCCESS 200 {object} app.Response{data=models.PlantClass} "获取成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "分类不存在"
// @Router /plant-classes/{id} [get]
func (h *Handler) GetPlantClass(c *gin.Context) {
	classIDStr := c.Param("id")
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的分类ID")
		return
	}

	class, err := h.PlantClassService.GetByID(classID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取分类信息失败: "+err.Error())
		return
	}

	if class == nil {
		app.Error(c, e.NOT_FOUND, "分类不存在")
		return
	}

	app.SUCCESS(c, class)
}

// GetPlantClasses 获取植物分类列表
// @Summary 获取植物分类列表
// @Description 分页获取所有植物分类信息
// @Tags 植物分类
// @Accept json
// @Produce json
// @Param page query int false "页码(默认1)" default(1)
// @Param size query int false "每页条数(默认10)" default(10)
// @Param with_count query bool false "是否包含植物数量(默认false)" default(false)
// @SUCCESS 200 {object} app.Response{data=app.PagedResult{list=[]models.PlantClass}} "获取成功"
// @SUCCESS 200 {object} app.Response{data=[]map[string]interface{}} "包含植物数量的分类列表"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Router /plant-classes [get]
func (h *Handler) GetPlantClasses(c *gin.Context) {
	// 检查是否需要包含植物数量
	withCount := c.Query("with_count") == "true"

	if withCount {
		// 返回包含植物数量的分类列表
		classes, err := h.PlantClassService.GetAllWithPlantCount()
		if err != nil {
			app.Error(c, e.INTERNAL_SERVER, "获取分类列表失败: "+err.Error())
			return
		}

		app.SUCCESS(c, classes)
		return
	}

	// 标准分页查询
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

	classes, total, err := h.PlantClassService.GetAll(page, pageSize)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取分类列表失败: "+err.Error())
		return
	}

	app.SUCCESS(c, app.PagedResult{
		List:      classes,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: (total + pageSize - 1) / pageSize,
	})
}

// CreatePlantClass 创建植物分类
// @Summary 创建植物分类
// @Description 创建新的植物分类信息
// @Tags 植物分类
// @Accept json
// @Produce json
// @Param class body models.PlantClass true "分类信息"
// @SUCCESS 200 {object} app.Response{data=models.PlantClass} "创建成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 500 {object} app.Response{data=string} "创建失败"
// @Security ApiKeyAuth
// @Router /plant-classes [post]
func (h *Handler) CreatePlantClass(c *gin.Context) {
	var class models.PlantClass

	if err := c.ShouldBindJSON(&class); err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的请求参数: "+err.Error())
		return
	}

	if class.Name == "" {
		app.Error(c, e.BAD_REQUEST, "分类名称不能为空")
		return
	}

	if err := h.PlantClassService.Create(&class); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "创建分类失败: "+err.Error())
		return
	}

	app.SUCCESS(c, class)
}

// UpdatePlantClass 更新植物分类信息
// @Summary 更新植物分类信息
// @Description 更新指定植物分类的信息
// @Tags 植物分类
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param class body models.PlantClass true "分类信息"
// @SUCCESS 200 {object} app.Response{data=models.PlantClass} "更新成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "分类不存在"
// @Failure 500 {object} app.Response{data=string} "更新失败"
// @Security ApiKeyAuth
// @Router /plant-classes/{id} [put]
func (h *Handler) UpdatePlantClass(c *gin.Context) {
	classIDStr := c.Param("id")
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的分类ID")
		return
	}

	// 检查分类是否存在
	existingClass, err := h.PlantClassService.GetByID(classID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取分类信息失败: "+err.Error())
		return
	}
	if existingClass == nil {
		app.Error(c, e.NOT_FOUND, "分类不存在")
		return
	}

	// 绑定请求体
	var class models.PlantClass
	if err := c.ShouldBindJSON(&class); err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的请求参数: "+err.Error())
		return
	}

	// 确保ID正确
	class.ID = classID

	if class.Name == "" {
		app.Error(c, e.BAD_REQUEST, "分类名称不能为空")
		return
	}

	if err := h.PlantClassService.Update(&class); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "更新分类失败: "+err.Error())
		return
	}

	app.SUCCESS(c, class)
}

// DeletePlantClass 删除植物分类
// @Summary 删除植物分类
// @Description 删除指定植物分类
// @Tags 植物分类
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @SUCCESS 200 {object} app.Response{data=map[string]interface{}} "删除成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "分类不存在"
// @Failure 500 {object} app.Response{data=string} "删除失败"
// @Security ApiKeyAuth
// @Router /plant-classes/{id} [delete]
func (h *Handler) DeletePlantClass(c *gin.Context) {
	classIDStr := c.Param("id")
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的分类ID")
		return
	}

	// 检查分类是否存在
	existingClass, err := h.PlantClassService.GetByID(classID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取分类信息失败: "+err.Error())
		return
	}
	if existingClass == nil {
		app.Error(c, e.NOT_FOUND, "分类不存在")
		return
	}

	if err := h.PlantClassService.Delete(classID); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "删除分类失败: "+err.Error())
		return
	}

	app.SUCCESS(c, gin.H{"message": "删除成功"})
}

// UploadPlantClassImage 上传植物分类图片
// @Summary 上传植物分类图片
// @Description 为指定植物分类上传图片
// @Tags 植物分类,图片
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "分类ID"
// @Param image formData file true "图片文件"
// @SUCCESS 200 {object} app.Response{data=map[string]string} "上传成功"
// @Failure 400 {object} app.Response{data=string} "请求参数错误"
// @Failure 404 {object} app.Response{data=string} "分类不存在"
// @Failure 500 {object} app.Response{data=string} "上传失败"
// @Security ApiKeyAuth
// @Router /plant-classes/{id}/images [post]
func (h *Handler) UploadPlantClassImage(c *gin.Context) {
	classIDStr := c.Param("id")
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "无效的分类ID")
		return
	}

	// 获取上传文件
	file, err := c.FormFile("image")
	if err != nil {
		app.Error(c, e.BAD_REQUEST, "请选择要上传的图片")
		return
	}

	// 上传图片并获取URL
	url, err := h.PlantClassService.UploadClassImage(classID, file)
	if err != nil {
		app.Error(c, e.ERROR_UPLOAD_FAILED, "上传图片失败: "+err.Error())
		return
	}

	// 获取当前分类信息
	class, err := h.PlantClassService.GetByID(classID)
	if err != nil {
		app.Error(c, e.INTERNAL_SERVER, "获取分类信息失败: "+err.Error())
		return
	}

	// 处理多图URL列表
	var srcList []string
	if class.SrcList != "" {
		srcList = strings.Split(class.SrcList, ",")
	}
	srcList = append(srcList, url)
	class.SrcList = strings.Join(srcList, ",")

	// 更新分类信息
	if err := h.PlantClassService.Update(class); err != nil {
		app.Error(c, e.INTERNAL_SERVER, "更新分类图片信息失败: "+err.Error())
		return
	}

	app.SUCCESS(c, gin.H{
		"url": url,
	})
}
