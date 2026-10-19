package service

import (
	"fmt"
	"strings"

	"github.com/rcolejnr/jjk-api/internal/models"
	"github.com/rcolejnr/jjk-api/internal/repository"
)

// TerminologyService defines the contract for Terminology business logic.
type TerminologyService interface {
	Insert(*models.CreateTerminologyRequest) (int64, error)
	Update(int64, *models.UpdateTerminologyRequest) error
	FindAll() ([]*models.Terminology, error)
	FindByID(int64) (*models.Terminology, error)
}

type terminologyService struct {
	repo repository.TerminologyRepository
}

func NewTerminologyService(repo repository.TerminologyRepository) TerminologyService {
	return &terminologyService{repo: repo}
}

func (s *terminologyService) Insert(req *models.CreateTerminologyRequest) (int64, error) {
	if term := strings.TrimSpace(req.Term); term == "" {
		return 0, fmt.Errorf("validation: term can't be empty")
	}
	if meaning := strings.TrimSpace(req.Meaning); meaning == "" {
		return 0, fmt.Errorf("validation: meaning can't be empty")
	}

	return s.Insert(req)
}

func (s *terminologyService) Update(id int64, req *models.UpdateTerminologyRequest) error {

	if req.Term != nil {
		trimmed := strings.TrimSpace(*req.Term)
		if trimmed == "" {
			return fmt.Errorf("validation: term can't be empty")
		}
	}
	if req.Meaning != nil {
		trimmed := strings.TrimSpace(*req.Meaning)
		if trimmed == "" {
			return fmt.Errorf("validation: meaning can't be empty")
		}
	}

	if req.Term == nil && req.Meaning == nil {
		return fmt.Errorf("validation: all update fields can't be empty")
	}

	return s.Update(id, req)
}

func (s *terminologyService) FindAll() ([]*models.Terminology, error) {
	terminologies, err := s.FindAll()
	if err != nil {
		return nil, err
	}

	return terminologies, nil
}

func (s *terminologyService) FindByID(id int64) (*models.Terminology, error) {
	term, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if term == nil {
		return nil, fmt.Errorf("terminology_service.FindByID: term with id %d does not exist", id)
	}
	return term, nil
}
