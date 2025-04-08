package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"botanical-api2/internal/models"
	"botanical-api2/internal/repository"
	"botanical-api2/pkg/setting"
	"botanical-api2/pkg/utils"
)

// PlantVegetationService 植物信息服务接口
type PlantVegetationService interface {
	// 创建植物信息
	Create(plant *models.PlantVegetation) error

	// 获取单个植物信息
	GetByID(id int) (*models.PlantVegetation, error)

	// 获取所有植物信息（分页）
	GetAll(page, pageSize int) ([]models.PlantVegetation, int, error)

	// 按分类获取植物信息
	GetByClass(classID int, page, pageSize int) ([]models.PlantVegetation, int, error)

	// 更新植物信息
	Update(plant *models.PlantVegetation) error

	// 删除植物信息
	Delete(id int) error

	// 获取园区内植物信息
	GetByParkID(parkID int, page, pageSize int) ([]models.PlantVegetation, int, error)

	// 上传植物图片
	UploadPlantImage(plantID int, file *multipart.FileHeader) (string, error)
}

// PlantVegetationServiceImpl 植物信息服务实现
type PlantVegetationServiceImpl struct {
	Repositories *repository.Repositories
}

// NewPlantVegetationService 创建植物信息服务
func NewPlantVegetationService(repositories *repository.Repositories) PlantVegetationService {
	return &PlantVegetationServiceImpl{
		Repositories: repositories,
	}
}

// Create 创建植物信息
func (s *PlantVegetationServiceImpl) Create(plant *models.PlantVegetation) error {
	return s.Repositories.PlantVegetation.Create(plant)
}

// GetByID 获取单个植物信息
func (s *PlantVegetationServiceImpl) GetByID(id int) (*models.PlantVegetation, error) {
	return s.Repositories.PlantVegetation.GetByID(id)
}

// GetAll 获取所有植物信息（分页）
func (s *PlantVegetationServiceImpl) GetAll(page, pageSize int) ([]models.PlantVegetation, int, error) {
	return s.Repositories.PlantVegetation.GetAll(page, pageSize)
}

// GetByClass 按分类获取植物信息
func (s *PlantVegetationServiceImpl) GetByClass(classID int, page, pageSize int) ([]models.PlantVegetation, int, error) {
	return s.Repositories.PlantVegetation.GetByClass(classID, page, pageSize)
}

// Update 更新植物信息
func (s *PlantVegetationServiceImpl) Update(plant *models.PlantVegetation) error {
	return s.Repositories.PlantVegetation.Update(plant)
}

// Delete 删除植物信息
func (s *PlantVegetationServiceImpl) Delete(id int) error {
	return s.Repositories.PlantVegetation.Delete(id)
}

// GetByParkID 获取园区内植物信息
func (s *PlantVegetationServiceImpl) GetByParkID(parkID int, page, pageSize int) ([]models.PlantVegetation, int, error) {
	return s.Repositories.PlantVegetation.GetByParkID(parkID, page, pageSize)
}

// UploadPlantImage 上传植物图片
func (s *PlantVegetationServiceImpl) UploadPlantImage(plantID int, file *multipart.FileHeader) (string, error) {
	// 检查植物是否存在
	plant, err := s.Repositories.PlantVegetation.GetByID(plantID)
	if err != nil {
		return "", fmt.Errorf("获取植物信息失败: %w", err)
	}
	if plant == nil {
		return "", errors.New("植物不存在")
	}

	// 检查文件类型
	if !utils.IsImageFile(file.Filename) {
		return "", errors.New("只允许上传图片文件")
	}

	// 生成唯一文件名
	fileName := fmt.Sprintf("plant_%d_%d%s", plantID, time.Now().UnixNano(), filepath.Ext(file.Filename))
	filePath := filepath.Join(setting.UploadPath, "plants", fileName)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 保存文件
	if err := utils.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 构建URL
	url := fmt.Sprintf("%s/uploads/plants/%s", setting.ServerDomain, fileName)

	return url, nil
}
