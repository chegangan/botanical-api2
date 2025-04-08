package repository

import (
	"errors"

	"botanical-api2/internal/models"

	"github.com/jinzhu/gorm"
)

// PlantNoticeRepository 植物资讯存储库实现
type PlantNoticeRepository struct {
	DB *gorm.DB
}

// NewPlantNoticeRepository 创建植物资讯存储库
func NewPlantNoticeRepository(db *gorm.DB) *PlantNoticeRepository {
	return &PlantNoticeRepository{DB: db}
}

// Create 创建植物资讯
func (r *PlantNoticeRepository) Create(notice *models.PlantNotice) error {
	return r.DB.Create(notice).Error
}

// GetByID 根据ID获取植物资讯
func (r *PlantNoticeRepository) GetByID(id int) (*models.PlantNotice, error) {
	var notice models.PlantNotice
	if err := r.DB.Where("id = ?", id).First(&notice).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &notice, nil
}

// GetAll 获取所有植物资讯(分页)
func (r *PlantNoticeRepository) GetAll(page, pageSize int, isRecommend *int) ([]models.PlantNotice, int, error) {
	var notices []models.PlantNotice
	var total int
	query := r.DB.Model(&models.PlantNotice{})

	// 如果指定了推荐状态，则添加筛选条件
	if isRecommend != nil {
		query = query.Where("is_recommend = ?", *isRecommend)
	}

	// 获取总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 构建查询
	query = r.DB.Model(&models.PlantNotice{})
	if isRecommend != nil {
		query = query.Where("is_recommend = ?", *isRecommend)
	}

	// 分页查询
	if err := query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("skim_sort DESC, id DESC").
		Find(&notices).Error; err != nil {
		return nil, 0, err
	}

	return notices, total, nil
}

// Update 更新植物资讯
func (r *PlantNoticeRepository) Update(notice *models.PlantNotice) error {
	return r.DB.Save(notice).Error
}

// Delete 删除植物资讯
func (r *PlantNoticeRepository) Delete(id int) error {
	return r.DB.Delete(&models.PlantNotice{}, id).Error
}

// IncrementSkimSort 增加浏览量
func (r *PlantNoticeRepository) IncrementSkimSort(id int) error {
	return r.DB.Model(&models.PlantNotice{}).
		Where("id = ?", id).
		UpdateColumn("skim_sort", gorm.Expr("skim_sort + ?", 1)).
		Error
}

// ToggleRecommend 切换推荐状态
func (r *PlantNoticeRepository) ToggleRecommend(id int, isRecommend int) error {
	return r.DB.Model(&models.PlantNotice{}).
		Where("id = ?", id).
		UpdateColumn("is_recommend", isRecommend).
		Error
}
