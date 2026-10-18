package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/rcolejnr/jjk-api/models"
)

// CharacterRepository defines the contract for Character data access.
type CharacterRepository interface {
	Insert(*models.CreateCharacterRequest) (int64, error)
	Update(int64, *models.UpdateCharacterRequest) error
	FindAll() ([]*models.Character, error)
	FindByID(int64) (*models.Character, error)
	Delete(int64) error
}

type sqlCharacterRepository struct {
	db *sql.DB
	sync.RWMutex
}

func NewCharacterRepository(db *sql.DB) CharacterRepository {
	return &sqlCharacterRepository{
		db: db,
	}
}

func (r *sqlCharacterRepository) Insert(req *models.CreateCharacterRequest) (int64, error) {
	r.Lock()
	defer r.Unlock()
	sql, err := r.db.Prepare(`INSERT INTO character VALUES(NULL,?,?,?,?,?, ?,?,?,?,?, ?,?,?,?)`)
	if err != nil {
		return 0, fmt.Errorf("character_repo.Insert: %w", err)
	}
	defer sql.Close()

	result, err := sql.Exec(
		req.Name,
		req.Description,
		req.Image,
		req.Species,
		req.Status,
		req.Birthday,
		req.Age,
		req.Gender,
		req.Height,
		req.Weight,
		req.HairColor,
		req.EyeColor,
		req.Grade,
		req.Clan,
		req.FirstAppearance,
	)
	if err != nil {
		return 0, fmt.Errorf("character_repo.Insert: %w", err)
	}

	var id int64
	if id, err = result.LastInsertId(); err != nil {
		return 0, fmt.Errorf("character_repo.Insert: %w", err)
	}

	// "KIARA"

	return id, nil
}

func (r *sqlCharacterRepository) Update(id int64, req *models.UpdateCharacterRequest) error {
	r.Lock()
	defer r.Unlock()

	setClauses := []string{}
	args := []any{}

	if req.Name != nil {
		setClauses = append(setClauses, "name = ?")
		args = append(args, *req.Name)
	}

	if req.Description != nil {
		setClauses = append(setClauses, "description = ?")
		args = append(args, *req.Description)
	}

	if req.Image != nil {
		setClauses = append(setClauses, "image = ?")
		args = append(args, *req.Image)
	}

	if req.Species != nil {
		setClauses = append(setClauses, "species = ?")
		args = append(args, *req.Species)
	}

	if req.Status != nil {
		setClauses = append(setClauses, "status = ?")
		args = append(args, *req.Status)
	}

	if req.Birthday != nil {
		setClauses = append(setClauses, "birthday = ?")
		args = append(args, *req.Birthday)
	}

	if req.Age != nil {
		setClauses = append(setClauses, "age = ?")
		args = append(args, *req.Age)
	}

	if req.Gender != nil {
		setClauses = append(setClauses, "gender = ?")
		args = append(args, *req.Gender)
	}

	if req.Height != nil {
		setClauses = append(setClauses, "height = ?")
		args = append(args, *req.Height)
	}

	if req.Weight != nil {
		setClauses = append(setClauses, "weight = ?")
		args = append(args, *req.Weight)
	}

	if req.HairColor != nil {
		setClauses = append(setClauses, "hair_color = ?")
		args = append(args, *req.HairColor)
	}

	if req.EyeColor != nil {
		setClauses = append(setClauses, "eye_color = ?")
		args = append(args, *req.EyeColor)
	}

	if req.Grade != nil {
		setClauses = append(setClauses, "grade = ?")
		args = append(args, *req.Grade)
	}

	if req.Clan != nil {
		setClauses = append(setClauses, "clan = ?")
		args = append(args, *req.Clan)
	}

	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE character SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	args = append(args, id)

	sql, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("character_repo.Update: %w", err)
	}
	defer sql.Close()

	_, err = sql.Exec(args...)
	return err
}

func (r *sqlCharacterRepository) FindAll() ([]*models.Character, error) {
	r.RLock()
	defer r.RUnlock()

	query := `
		SELECT id, name, description, image, species, status, birthday, 
			age, gender, height, weight, hair_color, eye_color, grade, clan, first_appearance_ep,
		FROM character
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("character_repo.FindAll: %w", err)
	}
	defer rows.Close()

	var characters []*models.Character
	for rows.Next() {
		var c models.Character
		var ep int64

		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.Image,
			&c.Species,
			&c.Status,
			&c.Birthday,
			&c.Age,
			&c.Gender,
			&c.Height,
			&c.Weight,
			&c.HairColor,
			&c.EyeColor,
			&c.Grade,
			&c.Clan,
			&ep,
		); err != nil {
			return nil, fmt.Errorf("character_repo.FindByID: %w", err)
		}

		c.FirstAppearance = &models.FirstAppearance{EpisodeID: ep}

		characters = append(characters, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("character_repo.FindAll: %w", err)
	}

	for _, char := range characters {
		// Query Domain Expansion
		var domain models.DomainExpansion
		row := r.db.QueryRow("SELECT id, character_id, name, description FROM domain_expansion WHERE character_id = ?", char.ID)
		err = row.Scan(&domain.ID, &domain.CharacterID, &domain.Name, &domain.Description)
		if err == nil {
			char.Domain = &domain
		} else if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("character_repo.FindByID: %w", err)
		}

		ctRows, err := r.db.Query("SELECT id, character_id, name, description, type FROM cursed_techniques WHERE character_id = ?", char.ID)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("character_repo.FindAll: %w", err)
		}

		// Query Cursed Techniques
		func() {
			defer ctRows.Close()
			for ctRows.Next() {
				var ct models.CursedTechnique
				if err := ctRows.Scan(&ct.ID, &ct.CharacterID, &ct.Name, &ct.Description, &ct.Type); err != nil {
					return
				}
				char.CursedTechniques = append(char.CursedTechniques, ct)
			}
		}()

		// Query Episode
		var episode models.Episode
		if err := r.db.QueryRow("SELECT id, title FROM episode WHERE id = ?", char.FirstAppearance.EpisodeID).Scan(&episode.ID, &episode.Title); err != nil {
			return nil, fmt.Errorf("character_repo.FindAll: %w", err)
		}
		char.FirstAppearance = &models.FirstAppearance{EpisodeID: episode.ID, EpisodeTitle: episode.Title}
	}

	return characters, nil
}

func (r *sqlCharacterRepository) FindByID(id int64) (*models.Character, error) {
	r.RLock()
	defer r.RUnlock()

	c := &models.Character{}
	var ep int64
	row := r.db.QueryRow("SELECT * FROM character WHERE id = ?", id)
	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Description,
		&c.Image,
		&c.Species,
		&c.Status,
		&c.Birthday,
		&c.Age,
		&c.Gender,
		&c.Height,
		&c.Weight,
		&c.HairColor,
		&c.EyeColor,
		&c.Grade,
		&c.Clan,
		&ep,
	)
	if err != nil {
		return nil, fmt.Errorf("character_repo.FindByID: %w", err)
	}

	// Query Domain Expansion
	domain := models.DomainExpansion{}
	row = r.db.QueryRow("SELECT id, character_id, name, description FROM domain_expansion WHERE character_id = ?", id)
	err = row.Scan(&domain.ID, &domain.CharacterID, &domain.Name, &domain.Description)
	if err == nil {
		c.Domain = &domain
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("character_repo.FindByID: %w", err)
	}

	// Query Cursed Techniques
	rows, err := r.db.Query("SELECT id, character_id, name, description, type FROM cursed_techniques WHERE character_id = ?", id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("character_repo.FindByID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ct models.CursedTechnique
		err := rows.Scan(&ct.ID, &ct.CharacterID, &ct.Name, &ct.Description, &ct.Type)
		if err != nil {
			return nil, fmt.Errorf("character_repo.FindByID: %w", err)
		}
		c.CursedTechniques = append(c.CursedTechniques, ct)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("character_repo.FindByID: %w", err)
	}

	// Query Episode
	var episode models.Episode
	if err := r.db.QueryRow("SELECT id, title FROM episode WHERE id = ?", &ep).Scan(&episode.ID, &episode.Title); err != nil {
		return nil, fmt.Errorf("character_repo.FindByID: %w", err)
	}
	c.FirstAppearance = &models.FirstAppearance{EpisodeID: episode.ID, EpisodeTitle: episode.Title}

	return c, nil
}

func (r *sqlCharacterRepository) Delete(id int64) error {
	r.Lock()
	defer r.Unlock()

	sql := "DELETE FORM character WHERE id = ?"
	result, err := r.db.Exec(sql, id)

	if err != nil {
		return fmt.Errorf("character_repo.Delete: %w", err)
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("character_repo.Delete: no character found with id: %d", id)
	}

	return nil
}
