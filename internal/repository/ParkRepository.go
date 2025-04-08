package repository

import (
	"errors"

	"botanical-api2/internal/models"

	"github.com/jinzhu/gorm"
)

// ParkRepository 园区存储库实现
type ParkRepository struct {
	DB *gorm.DB
}

// NewParkRepository 创建园区存储库
func NewParkRepository(db *gorm.DB) *ParkRepository {
	return &ParkRepository{DB: db}
}

// Create 创建园区
func (r *ParkRepository) Create(park *models.Park) error {
	return r.DB.Create(park).Error
}

// GetByID 根据ID获取园区
func (r *ParkRepository) GetByID(id int) (*models.Park, error) {
	var park models.Park
	if err := r.DB.Where("id = ?", id).First(&park).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &park, nil
}

// GetAll 获取所有园区(分页)
func (r *ParkRepository) GetAll(page, pageSize int) ([]models.Park, int, error) {
	var parks []models.Park
	var total int

	// 获取总记录数
	if err := r.DB.Model(&models.Park{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := r.DB.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id ASC").
		Find(&parks).Error; err != nil {
		return nil, 0, err
	}

	return parks, total, nil
}

// Update 更新园区
func (r *ParkRepository) Update(park *models.Park) error {
	return r.DB.Save(park).Error
}

// Delete 删除园区
func (r *ParkRepository) Delete(id int) error {
	// 先删除关联的园区植物关系
	if err := r.DB.Exec("DELETE FROM plant_park_vegetation WHERE park_id = ?", id).Error; err != nil {
		return err
	}

	// 删除园区
	return r.DB.Delete(&models.Park{}, id).Error
}

// AddPlantToPark 将植物添加到园区
func (r *ParkRepository) AddPlantToPark(parkID, plantID int) error {
	// 检查关联是否已存在
	var count int
	if err := r.DB.Model(&models.PlantParkVegetation{}).
		Where("park_id = ? AND plant_id = ?", parkID, plantID).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("植物已添加到该园区")
	}

	// 添加关联
	relation := &models.PlantParkVegetation{
		ParkID:  parkID,
		PlantID: plantID,
	}
	return r.DB.Create(relation).Error
}

// RemovePlantFromPark 从园区移除植物
func (r *ParkRepository) RemovePlantFromPark(parkID, plantID int) error {
	return r.DB.Where("park_id = ? AND plant_id = ?", parkID, plantID).
		Delete(&models.PlantParkVegetation{}).Error
}
