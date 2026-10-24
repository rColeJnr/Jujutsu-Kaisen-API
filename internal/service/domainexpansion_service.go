package service

import (
	"fmt"
	"strings"

	"github.com/rcolejnr/jjk-api/internal/models"
	"github.com/rcolejnr/jjk-api/internal/repository"
)

// DomainExpansionService defines the contract for domain expansion business logic.
type DomainExpansionService interface {
	Insert(*models.CreateDomainExpansionRequest) (int64, error)
	Update(int64, *models.UpdateDomainExpansionRequest) error
	FindAll(*models.DomainExpansionFilter) ([]*models.DomainExpansion, error)
	FindByID(int64) (*models.DomainExpansion, error)
}

type domainExpansionService struct {
	cRepo repository.CharacterRepository
	repo  repository.DomainExpansionRepository
}

func NewDomainExpansionService(cRepo repository.CharacterRepository, repo repository.DomainExpansionRepository) DomainExpansionService {
	return &domainExpansionService{cRepo: cRepo, repo: repo}
}

func (s *domainExpansionService) Insert(req *models.CreateDomainExpansionRequest) (int64, error) {
	if req.CharacterID < 1 {
		return 0, fmt.Errorf("validation: invalid character id")
	}

	_, err := s.cRepo.FindByID(req.CharacterID)
	if err != nil {
		return 0, fmt.Errorf("validation: character with id %d doesn't exist.", req.CharacterID)
	}

	if name := strings.TrimSpace(req.Name); name == "" {
		return 0, fmt.Errorf("validation: domain expansion name can't be empty")
	}
	if descp := strings.TrimSpace(req.Description); descp == "" {
		return 0, fmt.Errorf("validation: domain expansion description can't be empty")
	}

	return s.repo.Insert(req)
}

func (s *domainExpansionService) Update(id int64, req *models.UpdateDomainExpansionRequest) error {
	if *req.CharacterID < 1 {
		return fmt.Errorf("validation: invalid character id")
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			return fmt.Errorf("validation: domain expansion name can't be empty")
		}
	}
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		if trimmed == "" {
			return fmt.Errorf("validation: domain expansion description can't be empty")
		}
	}

	if req.Name == nil && req.Description == nil {
		return fmt.Errorf("validation: all update fields can't be empty§")
	}

	return s.repo.Update(id, req)
}

func (s *domainExpansionService) FindAll(f *models.DomainExpansionFilter) ([]*models.DomainExpansion, error) {
	domainsExpansion, err := s.repo.FindAll(f)
	if err != nil {
		return nil, err
	}

	return domainsExpansion, nil
}

func (s *domainExpansionService) FindByID(id int64) (*models.DomainExpansion, error) {
	domainExpansion, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if domainExpansion == nil {
		return nil, fmt.Errorf("dexpansion_service.FindByID: domain expansion with id %d does not exist", id)
	}
	return domainExpansion, nil
}
