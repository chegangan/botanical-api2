package repository

import (
	"errors"

	"botanical-api2/internal/models"

	"github.com/jinzhu/gorm"
)

// PlantVegetationRepository 植物信息存储库实现
type PlantVegetationRepository struct {
	DB *gorm.DB
}

// NewPlantVegetationRepository 创建植物信息存储库
func NewPlantVegetationRepository(db *gorm.DB) *PlantVegetationRepository {
	return &PlantVegetationRepository{DB: db}
}

// Create 创建植物信息
func (r *PlantVegetationRepository) Create(plant *models.PlantVegetation) error {
	return r.DB.Create(plant).Error
}

// GetByID 根据ID获取植物信息
func (r *PlantVegetationRepository) GetByID(id int) (*models.PlantVegetation, error) {
	var plant models.PlantVegetation
	if err := r.DB.Preload("Class").Where("id = ?", id).First(&plant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &plant, nil
}

// GetAll 获取所有植物信息(分页)
func (r *PlantVegetationRepository) GetAll(page, pageSize int) ([]models.PlantVegetation, int, error) {
	var plants []models.PlantVegetation
	var total int

	// 获取总记录数
	if err := r.DB.Model(&models.PlantVegetation{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := r.DB.Preload("Class").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("skim_sort DESC, id DESC").
		Find(&plants).Error; err != nil {
		return nil, 0, err
	}

	return plants, total, nil
}

// GetByClass 根据分类ID获取植物信息
func (r *PlantVegetationRepository) GetByClass(classID int, page, pageSize int) ([]models.PlantVegetation, int, error) {
	var plants []models.PlantVegetation
	var total int

	// 获取指定分类的总记录数
	if err := r.DB.Model(&models.PlantVegetation{}).Where("class_id = ?", classID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := r.DB.Preload("Class").
		Where("class_id = ?", classID).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("skim_sort DESC, id DESC").
		Find(&plants).Error; err != nil {
		return nil, 0, err
	}

	return plants, total, nil
}

// Update 更新植物信息
func (r *PlantVegetationRepository) Update(plant *models.PlantVegetation) error {
	return r.DB.Save(plant).Error
}

// Delete 删除植物信息
func (r *PlantVegetationRepository) Delete(id int) error {
	return r.DB.Delete(&models.PlantVegetation{}, id).Error
}

// GetByParkID 获取指定园区的植物信息
func (r *PlantVegetationRepository) GetByParkID(parkID int, page, pageSize int) ([]models.PlantVegetation, int, error) {
	var plants []models.PlantVegetation
	var total int

	// 通过关联表查询总数
	if err := r.DB.Model(&models.PlantParkVegetation{}).
		Where("park_id = ?", parkID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 通过关联表查询植物信息
	var plantIDs []int
	if err := r.DB.Model(&models.PlantParkVegetation{}).
		Where("park_id = ?", parkID).
		Pluck("plant_id", &plantIDs).Error; err != nil {
		return nil, 0, err
	}

	if len(plantIDs) == 0 {
		return []models.PlantVegetation{}, 0, nil
	}

	// 分页查询植物详细信息
	if err := r.DB.Preload("Class").
		Where("id IN (?)", plantIDs).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("skim_sort DESC, id DESC").
		Find(&plants).Error; err != nil {
		return nil, 0, err
	}

	return plants, total, nil
}
