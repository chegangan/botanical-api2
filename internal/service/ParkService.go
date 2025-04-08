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

// ParkService 园区服务接口
type ParkService interface {
	// 创建园区
	Create(park *models.Park) error

	// 获取单个园区
	GetByID(id int) (*models.Park, error)

	// 获取所有园区
	GetAll(page, pageSize int) ([]models.Park, int, error)

	// 更新园区
	Update(park *models.Park) error

	// 删除园区
	Delete(id int) error

	// 上传园区图片
	UploadParkImage(parkID int, file *multipart.FileHeader) (string, error)

	// 将植物添加到园区
	AddPlantToPark(parkID, plantID int) error

	// 从园区移除植物
	RemovePlantFromPark(parkID, plantID int) error
}

// ParkServiceImpl 园区服务实现
type ParkServiceImpl struct {
	Repositories *repository.Repositories
}

// NewParkService 创建园区服务
func NewParkService(repositories *repository.Repositories) ParkService {
	return &ParkServiceImpl{
		Repositories: repositories,
	}
}

// Create 创建园区
func (s *ParkServiceImpl) Create(park *models.Park) error {
	return s.Repositories.Park.Create(park)
}

// GetByID 获取单个园区
func (s *ParkServiceImpl) GetByID(id int) (*models.Park, error) {
	return s.Repositories.Park.GetByID(id)
}

// GetAll 获取所有园区（分页）
func (s *ParkServiceImpl) GetAll(page, pageSize int) ([]models.Park, int, error) {
	return s.Repositories.Park.GetAll(page, pageSize)
}

// Update 更新园区
func (s *ParkServiceImpl) Update(park *models.Park) error {
	return s.Repositories.Park.Update(park)
}

// Delete 删除园区
func (s *ParkServiceImpl) Delete(id int) error {
	return s.Repositories.Park.Delete(id)
}

// UploadParkImage 上传园区图片
func (s *ParkServiceImpl) UploadParkImage(parkID int, file *multipart.FileHeader) (string, error) {
	// 检查园区是否存在
	park, err := s.Repositories.Park.GetByID(parkID)
	if err != nil {
		return "", fmt.Errorf("获取园区信息失败: %w", err)
	}
	if park == nil {
		return "", errors.New("园区不存在")
	}

	// 检查文件类型
	if !utils.IsImageFile(file.Filename) {
		return "", errors.New("只允许上传图片文件")
	}

	// 生成唯一文件名
	fileName := fmt.Sprintf("park_%d_%d%s", parkID, time.Now().UnixNano(), filepath.Ext(file.Filename))
	filePath := filepath.Join(setting.UploadPath, "parks", fileName)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 保存文件
	if err := utils.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 构建URL
	url := fmt.Sprintf("%s/uploads/parks/%s", setting.ServerDomain, fileName)

	return url, nil
}

// AddPlantToPark 将植物添加到园区
func (s *ParkServiceImpl) AddPlantToPark(parkID, plantID int) error {
	// 检查园区是否存在
	park, err := s.Repositories.Park.GetByID(parkID)
	if err != nil {
		return fmt.Errorf("获取园区信息失败: %w", err)
	}
	if park == nil {
		return errors.New("园区不存在")
	}

	// 检查植物是否存在
	plant, err := s.Repositories.PlantVegetation.GetByID(plantID)
	if err != nil {
		return fmt.Errorf("获取植物信息失败: %w", err)
	}
	if plant == nil {
		return errors.New("植物不存在")
	}

	return s.Repositories.Park.AddPlantToPark(parkID, plantID)
}

// RemovePlantFromPark 从园区移除植物
func (s *ParkServiceImpl) RemovePlantFromPark(parkID, plantID int) error {
	return s.Repositories.Park.RemovePlantFromPark(parkID, plantID)
}
