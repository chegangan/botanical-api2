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

// PlantClassService 植物分类服务接口
type PlantClassService interface {
	// 创建植物分类
	Create(class *models.PlantClass) error

	// 获取单个植物分类
	GetByID(id int) (*models.PlantClass, error)

	// 获取所有植物分类
	GetAll(page, pageSize int) ([]models.PlantClass, int, error)

	// 更新植物分类
	Update(class *models.PlantClass) error

	// 删除植物分类
	Delete(id int) error

	// 获取所有分类及其植物数量
	GetAllWithPlantCount() ([]map[string]interface{}, error)

	// 上传分类图片
	UploadClassImage(classID int, file *multipart.FileHeader) (string, error)
}

// PlantClassServiceImpl 植物分类服务实现
type PlantClassServiceImpl struct {
	Repositories *repository.Repositories
}

// NewPlantClassService 创建植物分类服务
func NewPlantClassService(repositories *repository.Repositories) PlantClassService {
	return &PlantClassServiceImpl{
		Repositories: repositories,
	}
}

// Create 创建植物分类
func (s *PlantClassServiceImpl) Create(class *models.PlantClass) error {
	return s.Repositories.PlantClass.Create(class)
}

// GetByID 获取单个植物分类
func (s *PlantClassServiceImpl) GetByID(id int) (*models.PlantClass, error) {
	return s.Repositories.PlantClass.GetByID(id)
}

// GetAll 获取所有植物分类（分页）
func (s *PlantClassServiceImpl) GetAll(page, pageSize int) ([]models.PlantClass, int, error) {
	return s.Repositories.PlantClass.GetAll(page, pageSize)
}

// Update 更新植物分类
func (s *PlantClassServiceImpl) Update(class *models.PlantClass) error {
	return s.Repositories.PlantClass.Update(class)
}

// Delete 删除植物分类
func (s *PlantClassServiceImpl) Delete(id int) error {
	return s.Repositories.PlantClass.Delete(id)
}

// GetAllWithPlantCount 获取所有分类及其植物数量
func (s *PlantClassServiceImpl) GetAllWithPlantCount() ([]map[string]interface{}, error) {
	return s.Repositories.PlantClass.GetAllWithPlantCount()
}

// UploadClassImage 上传分类图片
func (s *PlantClassServiceImpl) UploadClassImage(classID int, file *multipart.FileHeader) (string, error) {
	// 检查分类是否存在
	class, err := s.Repositories.PlantClass.GetByID(classID)
	if err != nil {
		return "", fmt.Errorf("获取分类信息失败: %w", err)
	}
	if class == nil {
		return "", errors.New("分类不存在")
	}

	// 检查文件类型
	if !utils.IsImageFile(file.Filename) {
		return "", errors.New("只允许上传图片文件")
	}

	// 生成唯一文件名
	fileName := fmt.Sprintf("class_%d_%d%s", classID, time.Now().UnixNano(), filepath.Ext(file.Filename))
	filePath := filepath.Join(setting.UploadPath, "classes", fileName)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 保存文件
	if err := utils.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 构建URL
	url := fmt.Sprintf("%s/uploads/classes/%s", setting.ServerDomain, fileName)

	return url, nil
}
