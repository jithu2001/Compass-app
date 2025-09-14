package services

import (
	"compass-backend/internal/models"
	"compass-backend/internal/repositories"
)

type SpecificationService interface {
	CreateSpecification(spec *models.ProjectSpecification) error
	GetProjectSpecifications(projectID uint64) ([]models.ProjectSpecification, error)
	GetLatestSpecification(projectID uint64) (*models.ProjectSpecification, error)
}

type specificationService struct {
	specRepo repositories.SpecificationRepository
}

func NewSpecificationService(specRepo repositories.SpecificationRepository) SpecificationService {
	return &specificationService{
		specRepo: specRepo,
	}
}

func (s *specificationService) CreateSpecification(spec *models.ProjectSpecification) error {
	// Version number is automatically handled in the model's BeforeCreate hook
	return s.specRepo.Create(spec)
}

func (s *specificationService) GetProjectSpecifications(projectID uint64) ([]models.ProjectSpecification, error) {
	return s.specRepo.FindByProjectID(projectID)
}

func (s *specificationService) GetLatestSpecification(projectID uint64) (*models.ProjectSpecification, error) {
	return s.specRepo.FindLatestByProjectID(projectID)
}