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

// PlantNoticeService 植物资讯服务接口
type PlantNoticeService interface {
	// 创建植物资讯
	Create(notice *models.PlantNotice) error

	// 获取单个植物资讯
	GetByID(id int) (*models.PlantNotice, error)

	// 获取所有植物资讯（分页）
	GetAll(page, pageSize int, isRecommend *int) ([]models.PlantNotice, int, error)

	// 更新植物资讯
	Update(notice *models.PlantNotice) error

	// 删除植物资讯
	Delete(id int) error

	// 增加浏览量
	IncrementSkimSort(id int) error

	// 切换推荐状态
	ToggleRecommend(id int, isRecommend int) error

	// 上传资讯图片
	UploadNoticeImage(noticeID int, file *multipart.FileHeader) (string, error)
}

// PlantNoticeServiceImpl 植物资讯服务实现
type PlantNoticeServiceImpl struct {
	Repositories *repository.Repositories
}

// NewPlantNoticeService 创建植物资讯服务
func NewPlantNoticeService(repositories *repository.Repositories) PlantNoticeService {
	return &PlantNoticeServiceImpl{
		Repositories: repositories,
	}
}

// Create 创建植物资讯
func (s *PlantNoticeServiceImpl) Create(notice *models.PlantNotice) error {
	notice.CreatedAt = time.Now()
	notice.UpdatedAt = time.Now()
	return s.Repositories.PlantNotice.Create(notice)
}

// GetByID 获取单个植物资讯
func (s *PlantNoticeServiceImpl) GetByID(id int) (*models.PlantNotice, error) {
	return s.Repositories.PlantNotice.GetByID(id)
}

// GetAll 获取所有植物资讯（分页）
func (s *PlantNoticeServiceImpl) GetAll(page, pageSize int, isRecommend *int) ([]models.PlantNotice, int, error) {
	return s.Repositories.PlantNotice.GetAll(page, pageSize, isRecommend)
}

// Update 更新植物资讯
func (s *PlantNoticeServiceImpl) Update(notice *models.PlantNotice) error {
	// 获取现有记录
	existingNotice, err := s.GetByID(notice.ID)
	if err != nil {
		return err
	}
	if existingNotice == nil {
		return errors.New("资讯不存在")
	}

	// 保留原始创建时间
	notice.CreatedAt = existingNotice.CreatedAt
	notice.UpdatedAt = time.Now()
	return s.Repositories.PlantNotice.Update(notice)
}

// Delete 删除植物资讯
func (s *PlantNoticeServiceImpl) Delete(id int) error {
	return s.Repositories.PlantNotice.Delete(id)
}

// IncrementSkimSort 增加浏览量
func (s *PlantNoticeServiceImpl) IncrementSkimSort(id int) error {
	return s.Repositories.PlantNotice.IncrementSkimSort(id)
}

// ToggleRecommend 切换推荐状态
func (s *PlantNoticeServiceImpl) ToggleRecommend(id int, isRecommend int) error {
	if isRecommend != 0 && isRecommend != 1 {
		return errors.New("推荐状态只能为0或1")
	}
	return s.Repositories.PlantNotice.ToggleRecommend(id, isRecommend)
}

// UploadNoticeImage 上传资讯图片
func (s *PlantNoticeServiceImpl) UploadNoticeImage(noticeID int, file *multipart.FileHeader) (string, error) {
	// 检查资讯是否存在
	notice, err := s.Repositories.PlantNotice.GetByID(noticeID)
	if err != nil {
		return "", fmt.Errorf("获取资讯信息失败: %w", err)
	}
	if notice == nil {
		return "", errors.New("资讯不存在")
	}

	// 检查文件类型
	if !utils.IsImageFile(file.Filename) {
		return "", errors.New("只允许上传图片文件")
	}

	// 生成唯一文件名
	fileName := fmt.Sprintf("notice_%d_%d%s", noticeID, time.Now().UnixNano(), filepath.Ext(file.Filename))
	filePath := filepath.Join(setting.UploadPath, "notices", fileName)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 保存文件
	if err := utils.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 构建URL
	url := fmt.Sprintf("%s/uploads/notices/%s", setting.ServerDomain, fileName)

	return url, nil
}
