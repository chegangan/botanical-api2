package repository

import (
	"errors"

	"botanical-api2/internal/models"

	"github.com/jinzhu/gorm"
)

// PlantClassRepository 植物分类存储库实现
type PlantClassRepository struct {
	DB *gorm.DB
}

// NewPlantClassRepository 创建植物分类存储库
func NewPlantClassRepository(db *gorm.DB) *PlantClassRepository {
	return &PlantClassRepository{DB: db}
}

// Create 创建植物分类
func (r *PlantClassRepository) Create(class *models.PlantClass) error {
	return r.DB.Create(class).Error
}

// GetByID 根据ID获取植物分类
func (r *PlantClassRepository) GetByID(id int) (*models.PlantClass, error) {
	var class models.PlantClass
	if err := r.DB.Where("id = ?", id).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &class, nil
}

// GetAll 获取所有植物分类(分页)
func (r *PlantClassRepository) GetAll(page, pageSize int) ([]models.PlantClass, int, error) {
	var classes []models.PlantClass
	var total int

	// 获取总记录数
	if err := r.DB.Model(&models.PlantClass{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := r.DB.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id ASC").
		Find(&classes).Error; err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}

// Update 更新植物分类
func (r *PlantClassRepository) Update(class *models.PlantClass) error {
	return r.DB.Save(class).Error
}

// Delete 删除植物分类
func (r *PlantClassRepository) Delete(id int) error {
	// 检查是否有关联的植物
	var count int
	if err := r.DB.Model(&models.PlantVegetation{}).Where("class_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该分类下存在植物，无法删除")
	}

	return r.DB.Delete(&models.PlantClass{}, id).Error
}

// GetAllWithPlantCount 获取所有分类及其植物数量
func (r *PlantClassRepository) GetAllWithPlantCount() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 使用原生SQL查询
	rows, err := r.DB.Raw(`
        SELECT pc.id, pc.name, pc.src_list, COUNT(pv.id) as plant_count 
        FROM plant_class pc 
        LEFT JOIN plant_vegetation pv ON pc.id = pv.class_id 
        GROUP BY pc.id, pc.name, pc.src_list
        ORDER BY pc.id ASC
    `).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, srcList string
		var plantCount int
		if err := rows.Scan(&id, &name, &srcList, &plantCount); err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"id":          id,
			"name":        name,
			"src_list":    srcList,
			"plant_count": plantCount,
		})
	}

	return results, nil
}
