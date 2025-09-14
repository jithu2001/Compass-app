package repositories

import (
	"compass-backend/internal/models"
	"gorm.io/gorm"
)

type SpecificationRepository interface {
	Create(spec *models.ProjectSpecification) error
	FindByProjectID(projectID uint64) ([]models.ProjectSpecification, error)
	FindLatestByProjectID(projectID uint64) (*models.ProjectSpecification, error)
}

type specificationRepository struct {
	db *gorm.DB
}

func NewSpecificationRepository(db *gorm.DB) SpecificationRepository {
	return &specificationRepository{db: db}
}

func (r *specificationRepository) Create(spec *models.ProjectSpecification) error {
	return r.db.Create(spec).Error
}

func (r *specificationRepository) FindByProjectID(projectID uint64) ([]models.ProjectSpecification, error) {
	var specs []models.ProjectSpecification
	err := r.db.Where("project_id = ?", projectID).
		Preload("Creator").
		Order("version_no DESC").
		Find(&specs).Error
	return specs, err
}

func (r *specificationRepository) FindLatestByProjectID(projectID uint64) (*models.ProjectSpecification, error) {
	var spec models.ProjectSpecification
	err := r.db.Where("project_id = ?", projectID).
		Preload("Creator").
		Order("version_no DESC").
		First(&spec).Error
	if err != nil {
		return nil, err
	}
	return &spec, nil
}