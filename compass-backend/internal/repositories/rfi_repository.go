package repositories

import (
	"compass-backend/internal/models"
	"gorm.io/gorm"
)

type RFIRepository interface {
	Create(rfi *models.ProjectRFI) error
	FindByID(id uint64) (*models.ProjectRFI, error)
	FindByProjectID(projectID uint64) ([]models.ProjectRFI, error)
	Answer(id uint64, answer models.AnswerValue, answeredBy uint64) error
}

type rfiRepository struct {
	db *gorm.DB
}

func NewRFIRepository(db *gorm.DB) RFIRepository {
	return &rfiRepository{db: db}
}

func (r *rfiRepository) Create(rfi *models.ProjectRFI) error {
	return r.db.Create(rfi).Error
}

func (r *rfiRepository) FindByID(id uint64) (*models.ProjectRFI, error) {
	var rfi models.ProjectRFI
	err := r.db.Preload("Answerer").First(&rfi, id).Error
	if err != nil {
		return nil, err
	}
	return &rfi, nil
}

func (r *rfiRepository) FindByProjectID(projectID uint64) ([]models.ProjectRFI, error) {
	var rfis []models.ProjectRFI
	err := r.db.Where("project_id = ?", projectID).
		Preload("Answerer").
		Order("created_at DESC").
		Find(&rfis).Error
	return rfis, err
}

func (r *rfiRepository) Answer(id uint64, answer models.AnswerValue, answeredBy uint64) error {
	return r.db.Model(&models.ProjectRFI{}).Where("rfi_id = ?", id).Updates(map[string]interface{}{
		"answer_value": answer,
		"answered_by":  answeredBy,
	}).Error
}