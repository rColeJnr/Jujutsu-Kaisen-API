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
	FindAll(*models.EpisodeFilter) ([]*models.Episode, int, error)
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

	if err := validateTexts(
		text("title", req.Title, maxNameLen),
		imageURL("image", req.Image),
		text("air_date", req.AirDate, maxShortLen),
	); err != nil {
		return 0, err
	}

	return s.repo.Insert(req)
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

	if err := validateTexts(
		text("title", str(req.Title), maxNameLen),
		imageURL("image", str(req.Image)),
		text("air_date", str(req.AirDate), maxShortLen),
	); err != nil {
		return err
	}

	return s.repo.Update(id, req)
}

func (s *episodeService) FindAll(f *models.EpisodeFilter) ([]*models.Episode, int, error) {
	return s.repo.FindAll(f)
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
