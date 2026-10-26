package service

import (
	"fmt"
	"strings"

	"github.com/rcolejnr/jjk-api/internal/models"
	"github.com/rcolejnr/jjk-api/internal/repository"
)

// CursedTechniqueService defines the contract for Cursed Technique business logic.
type CursedTechniqueService interface {
	Insert(*models.CreateCursedTechniqueRequest) (int64, error)
	Update(int64, *models.UpdateCursedTechniqueRequest) error
	FindAll(*models.CursedTechniqueFilter) ([]*models.CursedTechnique, int, error)
	FindByID(int64) (*models.CursedTechnique, error)
}

type ctechniqueService struct {
	cRepo repository.CharacterRepository
	repo  repository.CursedTechniqueRepository
}

func NewCursedTechniqueService(cRepo repository.CharacterRepository, repo repository.CursedTechniqueRepository) CursedTechniqueService {
	return &ctechniqueService{cRepo: cRepo, repo: repo}
}

func (s *ctechniqueService) Insert(req *models.CreateCursedTechniqueRequest) (int64, error) {
	if req.CharacterID < 0 {
		return 0, fmt.Errorf("validation: invalid character id")
	}

	_, err := s.cRepo.FindByID(req.CharacterID)
	if err != nil {
		return 0, fmt.Errorf("validation: character with id %d doesn't exist.", req.CharacterID)
	}

	if name := strings.TrimSpace(req.Name); name == "" {
		return 0, fmt.Errorf("validation: technique name can't be empty")
	}
	if descp := strings.TrimSpace(req.Description); descp == "" {
		return 0, fmt.Errorf("validation: technique description can't be empty")
	}

	return s.repo.Insert(req)
}

func (s *ctechniqueService) Update(id int64, req *models.UpdateCursedTechniqueRequest) error {

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			return fmt.Errorf("validation: technique name can't be empty")
		}
	}
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		if trimmed == "" {
			return fmt.Errorf("validation: technique description can't be empty")
		}
	}

	if req.Name == nil && req.Description == nil && req.Type == nil {
		return fmt.Errorf("validation: all update fields can't be empty§")
	}

	return s.repo.Update(id, req)
}

func (s *ctechniqueService) FindAll(f *models.CursedTechniqueFilter) ([]*models.CursedTechnique, int, error) {
	return s.repo.FindAll(f)
}

func (s *ctechniqueService) FindByID(id int64) (*models.CursedTechnique, error) {
	cTechnique, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if cTechnique == nil {
		return nil, fmt.Errorf("ctechnique_service.FindByID: technique with id %d does not exist", id)
	}
	return cTechnique, nil
}
