package services

import (
	"errors"

	"compass-backend/internal/models"
	"compass-backend/internal/repositories"
	"compass-backend/db"
)

type ProjectService interface {
	CreateProject(project *models.Project) error
	CreateProjectWithDetails(project *models.Project, specifications []models.ProjectSpecification, rfis []models.ProjectRFI) error
	GetProject(projectID uint64) (*models.Project, error)
	ListProjects() ([]models.Project, error)
	UpdateProjectStatus(projectID uint64, status models.ProjectStatus, updatedBy uint64) error
	DeleteProject(projectID uint64, userRole models.UserRole) error
}

type projectService struct {
	projectRepo       repositories.ProjectRepository
	specificationRepo repositories.SpecificationRepository
	rfiRepo          repositories.RFIRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository, specRepo repositories.SpecificationRepository, rfiRepo repositories.RFIRepository) ProjectService {
	return &projectService{
		projectRepo:       projectRepo,
		specificationRepo: specRepo,
		rfiRepo:          rfiRepo,
	}
}

func (s *projectService) CreateProject(project *models.Project) error {
	return s.projectRepo.Create(project)
}

func (s *projectService) CreateProjectWithDetails(project *models.Project, specifications []models.ProjectSpecification, rfis []models.ProjectRFI) error {
	// Start a database transaction
	tx := db.GetDB().Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create the project first
	if err := tx.Create(project).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create specifications if provided
	for i := range specifications {
		specifications[i].ProjectID = project.ProjectID
		if err := tx.Create(&specifications[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Create RFIs if provided
	for i := range rfis {
		rfis[i].ProjectID = project.ProjectID
		if err := tx.Create(&rfis[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (s *projectService) GetProject(projectID uint64) (*models.Project, error) {
	return s.projectRepo.FindByID(projectID)
}

func (s *projectService) ListProjects() ([]models.Project, error) {
	return s.projectRepo.List()
}

func (s *projectService) UpdateProjectStatus(projectID uint64, status models.ProjectStatus, updatedBy uint64) error {
	// Validate status
	switch status {
	case models.StatusNotYetStarted, models.StatusProgress, models.StatusCompleted:
		// Valid status
	default:
		return errors.New("invalid project status")
	}

	return s.projectRepo.UpdateStatus(projectID, status, updatedBy)
}

func (s *projectService) DeleteProject(projectID uint64, userRole models.UserRole) error {
	if userRole != models.RoleAdmin {
		return errors.New("only admins can delete projects")
	}

	return s.projectRepo.Delete(projectID)
}