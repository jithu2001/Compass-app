package services

import (
	"errors"

	"compass-backend/internal/models"
	"compass-backend/internal/repositories"
)

type RFIService interface {
	CreateRFI(rfi *models.ProjectRFI) error
	GetRFI(rfiID uint64) (*models.ProjectRFI, error)
	GetProjectRFIs(projectID uint64) ([]models.ProjectRFI, error)
	AnswerRFI(rfiID uint64, answer models.AnswerValue, answeredBy uint64) error
}

type rfiService struct {
	rfiRepo repositories.RFIRepository
}

func NewRFIService(rfiRepo repositories.RFIRepository) RFIService {
	return &rfiService{
		rfiRepo: rfiRepo,
	}
}

func (s *rfiService) CreateRFI(rfi *models.ProjectRFI) error {
	return s.rfiRepo.Create(rfi)
}

func (s *rfiService) GetRFI(rfiID uint64) (*models.ProjectRFI, error) {
	return s.rfiRepo.FindByID(rfiID)
}

func (s *rfiService) GetProjectRFIs(projectID uint64) ([]models.ProjectRFI, error) {
	return s.rfiRepo.FindByProjectID(projectID)
}

func (s *rfiService) AnswerRFI(rfiID uint64, answer models.AnswerValue, answeredBy uint64) error {
	// Validate answer
	if answer != models.AnswerYes && answer != models.AnswerNo {
		return errors.New("invalid answer value")
	}

	// Check if RFI exists
	rfi, err := s.rfiRepo.FindByID(rfiID)
	if err != nil {
		return errors.New("RFI not found")
	}

	// Check if already answered
	if rfi.AnswerValue != nil {
		return errors.New("RFI already answered")
	}

	return s.rfiRepo.Answer(rfiID, answer, answeredBy)
}