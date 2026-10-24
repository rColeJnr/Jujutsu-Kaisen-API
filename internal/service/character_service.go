package service

import (
	"fmt"
	"strings"

	"github.com/rcolejnr/jjk-api/internal/models"
	"github.com/rcolejnr/jjk-api/internal/repository"
)

// CharacterService defines the contract for character business logic.
type CharacterService interface {
	Insert(*models.CreateCharacterRequest) (int64, error)
	Update(int64, *models.UpdateCharacterRequest) error
	FindAll() ([]*models.Character, error)
	FindByID(int64) (*models.Character, error)
	Delete(int64) error
}

type characterService struct {
	characterRepo repository.CharacterRepository
}

func NewCharacterService(repo repository.CharacterRepository) CharacterService {
	return &characterService{characterRepo: repo}
}

func (s *characterService) Insert(req *models.CreateCharacterRequest) (int64, error) {
	err := reqValidation(req)
	if err != nil {
		return 0, err
	}

	return s.characterRepo.Insert(req)
}

func (s *characterService) Update(id int64, req *models.UpdateCharacterRequest) error {
	if err := reqUpdateValidation(req); err != nil {
		return err
	}

	return s.characterRepo.Update(id, req)
}

func (s *characterService) FindAll() ([]*models.Character, error) {
	characters, err := s.characterRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return characters, nil
}

func (s *characterService) FindByID(id int64) (*models.Character, error) {
	character, err := s.characterRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if character == nil {
		return nil, fmt.Errorf("character_service.FindByID: character with id %d does not exist", id)
	}
	return character, nil
}

func (s *characterService) Delete(id int64) error {

	return s.characterRepo.Delete(id)
}

func reqValidation(req *models.CreateCharacterRequest) error {
	if trimmAndLowerCase(req.Name) == "" {
		return fmt.Errorf("validation: name can't be empty.")
	}

	if species := trimmAndLowerCase(req.Species); species == "" {
		return fmt.Errorf("validation: species can't be empty.")
	} else if !(species == "human" || species == "curse") {
		return fmt.Errorf("validation: invalid species")
	}

	if status := trimmAndLowerCase(req.Status); status == "" {
		return fmt.Errorf("validation: status can't be empty.")
	} else if !(status == "alive" || status == "dead" || status == "unknown") {
		return fmt.Errorf("validation: invalid status")
	}

	if gender := trimmAndLowerCase(req.Gender); !(gender == "male" || gender == "female") {
		return fmt.Errorf("validation: invalid gender")
	}

	if clan := trimmAndLowerCase(req.Clan); !(clan == "gojo" || clan == "zenin") {
		return fmt.Errorf("validation: invalid clan")
	}

	if grade := trimmAndLowerCase(req.Grade); !(grade == "special grade" || grade == "1st grade" || grade == "2nd grade" || grade == "3rd grade") {
		return fmt.Errorf("validation: invalid grade")
	}

	if req.Age < 0 {
		return fmt.Errorf("validation: age can't be less than 0.")
	}

	if req.Height <= 10 {
		return fmt.Errorf("validation: height can't be less than 10.")
	}

	if req.Weight <= 0 {
		return fmt.Errorf("validation: weight can't be less than 0.")
	}

	if req.FirstAppearanceEp <= 0 {
		return fmt.Errorf("validation: episode can't be less than 0, got: %d, %s.", req.FirstAppearanceEp, req.Name)
	}

	return nil
}

func reqUpdateValidation(req *models.UpdateCharacterRequest) error {
	if req.Name != nil {
		trimmed := trimmAndLowerCase(*req.Name)
		if trimmed == "" {
			return fmt.Errorf("validation: name can't be empty.")
		}
		req.Name = &trimmed
	}

	if req.Species != nil {
		trimmed := trimmAndLowerCase(*req.Species)
		if trimmed == "" {
			return fmt.Errorf("validation: species can't be empty.")
		}
		if trimmed != "human" && trimmed != "curse" {
			return fmt.Errorf("validation: invalid species")
		}
		req.Species = &trimmed
	}

	if req.Status != nil {
		trimmed := trimmAndLowerCase(*req.Status)
		if trimmed == "" {
			return fmt.Errorf("validation: status can't be empty.")
		}
		if trimmed != "alive" && trimmed != "dead" && trimmed != "unknown" {
			return fmt.Errorf("validation: invalid status")
		}
		req.Status = &trimmed
	}

	if req.Gender != nil {
		trimmed := trimmAndLowerCase(*req.Gender)
		if trimmed != "male" && trimmed != "femaile" {
			return fmt.Errorf("validation: invalid gender")
		}
		req.Gender = &trimmed
	}

	if req.Clan != nil {
		trimmed := trimmAndLowerCase(*req.Clan)
		if !(trimmed == "gojo" || trimmed == "zenin") {
			return fmt.Errorf("validation: invalid clan")
		}
		req.Status = &trimmed
	}

	if req.Grade != nil {
		trimmed := trimmAndLowerCase(*req.Grade)
		if !(trimmed == "special grade" || trimmed == "1st grade" || trimmed == "2nd grade" || trimmed == "3rd grade") {
			return fmt.Errorf("validation: invalid grade")
		}
		req.Grade = &trimmed
	}

	if req.Age != nil && *req.Age < 0 {
		return fmt.Errorf("validation: age can't be less than 0.")
	}

	if req.Height != nil && *req.Height <= 10 {
		return fmt.Errorf("validation: height can't be less than 10.")
	}

	if req.Weight != nil && *req.Weight <= 0 {
		return fmt.Errorf("validation: weight can't be less than 0.")
	}

	if req.Weight == nil && req.Height == nil && req.Age == nil && req.Grade == nil &&
		req.Clan == nil && req.Gender == nil && req.Status == nil && req.Species == nil && req.Name == nil {
		return fmt.Errorf("validation: no fields provided for update")
	}

	return nil
}

func trimmAndLowerCase(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
