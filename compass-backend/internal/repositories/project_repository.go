package repositories

import (
	"compass-backend/internal/models"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *models.Project) error
	FindByID(id uint64) (*models.Project, error)
	List() ([]models.Project, error)
	UpdateStatus(id uint64, status models.ProjectStatus, updatedBy uint64) error
	Delete(id uint64) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) FindByID(id uint64) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Creator").Preload("LastUpdater").First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) List() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Preload("Creator").Preload("LastUpdater").Find(&projects).Error
	return projects, err
}

func (r *projectRepository) UpdateStatus(id uint64, status models.ProjectStatus, updatedBy uint64) error {
	return r.db.Model(&models.Project{}).Where("project_id = ?", id).Updates(map[string]interface{}{
		"project_status": status,
		"last_updated_by": updatedBy,
	}).Error
}

func (r *projectRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Project{}, id).Error
}