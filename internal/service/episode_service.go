package service

import (
	"fmt"
	"strings"

	"github.com/rcolejnr/jjk-api/internal/models"
	"github.com/rcolejnr/jjk-api/internal/repository"
)

// CharacterService defines the contract for character business logic.
type EpisodeService interface {
	Insert(*models.CreateEpisodeRequest) (int64, error)
	Update(int64, *models.UpdateEpisodeRequest) error
	FindAll() ([]*models.Episode, error)
	FindByID(int64) (*models.Episode, error)
	Delete(int64) error
}

type episodeService struct {
	repo repository.EpisodeRepository
}

func NewEpiesodeService(repo repository.EpisodeRepository) EpisodeService {
	return &episodeService{repo: repo}
}

func (s *episodeService) Insert(req *models.CreateEpisodeRequest) (int64, error) {
	if req.Season < 1 {
		return 0, fmt.Errorf("validation: season can't be less then 1")
	}

	if title := strings.TrimSpace(req.Title); title == "" {
		return 0, fmt.Errorf("validation: title can't be empty")
	}

	if airDate := strings.TrimSpace(req.AirDate); airDate == "" {
		return 0, fmt.Errorf("validation: air date can't be empty")
	}

	if len(req.Characters) == 0 {
		return 0, fmt.Errorf("validation: characters can't be empty")
	}

	return s.Insert(req)
}

func (s *episodeService) Update(id int64, req *models.UpdateEpisodeRequest) error {
	if req.Season != nil && *req.Season < 1 {
		return fmt.Errorf("validation: season can't be less then 1")
	}

	if req.Title != nil {
		if title := strings.TrimSpace(*req.Title); title == "" {
			return fmt.Errorf("validation: title can't be empty")
		}
	}

	if req.AirDate != nil {
		if airDate := strings.TrimSpace(*req.AirDate); airDate == "" {
			return fmt.Errorf("validation: air date can't be empty")
		}
	}

	if req.Title == nil && req.Season == nil && req.AirDate == nil {
		return fmt.Errorf("validation: no fields provided for update")
	}

	return s.Update(id, req)
}

func (s *episodeService) FindAll() ([]*models.Episode, error) {
	episodes, err := s.FindAll()
	if err != nil {
		return nil, err
	}

	return episodes, nil
}

func (s *episodeService) FindByID(id int64) (*models.Episode, error) {
	episode, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if episode == nil {
		return nil, fmt.Errorf("episode_service.FindByID: episode with id %d does not exist", id)
	}
	return episode, nil
}

func (s *episodeService) Delete(id int64) error {
	episode, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if episode == nil {
		return fmt.Errorf("episode_service.Delete: episode with id %d does not exist", id)
	}
	return s.repo.Delete(id)
}
