package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/rcolejnr/jjk-api/internal/models"
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

	res, err := sql.Exec(args...)

	_, err = res.RowsAffected()

	return err
}

func (r *sqlCharacterRepository) FindAll() ([]*models.Character, error) {
	r.RLock()
	defer r.RUnlock()

	query := `
		SELECT c.id, c.name, c.description, c.image, c.species, c.status, c.birthday, 
			c.age, c.gender, c.height, c.weight, c.hair_color, c.eye_color, c.grade, c.clan, c.first_appearance_ep,
			d.id, d.name, d.description,
			e.title
		FROM character c
		LEFT JOIN domain_expansion d ON c.id = d.character_id
		LEFT JOIN episode e ON c.first_appearance_ep = e.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("character_repo.FindAll: %w", err)
	}
	defer rows.Close()

	var characters []*models.Character
	for rows.Next() {
		var c models.Character
		var ep int64 // first appearance episode
		var e models.Episode
		var dID sql.NullInt64 // domain expansion
		var dName, dDescription sql.NullString

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
			&dID,
			&dName,
			&dDescription,
			&e.Title,
		); err != nil {
			return nil, fmt.Errorf("character_repo.FindAll: %w", err)
		}

		c.FirstAppearance = &models.FirstAppearance{EpisodeID: ep, EpisodeTitle: e.Title}

		if dID.Valid {
			c.Domain = &models.DomainExpansion{
				ID:          dID.Int64,
				CharacterID: c.ID,
				Name:        dName.String,
				Description: dDescription.String,
			}
		}

		c.CursedTechniques = []models.CursedTechnique{}

		characters = append(characters, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("character_repo.FindAll: %w", err)
	}

	if len(characters) == 0 {
		return characters, nil
	}

	// Query Cursed Techniques
	for _, char := range characters {
		ctRows, err := r.db.Query("SELECT id, character_id, name, description, type FROM cursed_techniques WHERE character_id = ?", char.ID)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("character_repo.FindAll: %w", err)
		}

		ctErr := func() error {
			defer ctRows.Close()
			for ctRows.Next() {
				var ct models.CursedTechnique
				if err := ctRows.Scan(&ct.ID, &ct.CharacterID, &ct.Name, &ct.Description, &ct.Type); err != nil {
					return err
				}
				char.CursedTechniques = append(char.CursedTechniques, ct)
			}

			return ctRows.Err()
		}()

		if ctErr != nil {
			return nil, fmt.Errorf("character_repo.FindAll: %w", ctErr)
		}
	}

	return characters, nil
}

func (r *sqlCharacterRepository) FindByID(id int64) (*models.Character, error) {
	r.RLock()
	defer r.RUnlock()

	query := `
		SELECT c.id, c.name, c.description, c.image, c.species, c.status, c.birthday, 
			c.age, c.gender, c.height, c.weight, c.hair_color, c.eye_color, c.grade, c.clan, c.first_appearance_ep,
			d.id, d.name, d.description,
			e.title,
			t.id, t.name, t.description, t.type
		FROM character c
		LEFT JOIN domain_expansion d ON c.id = d.character_id
		LEFT JOIN episode e ON c.first_appearance_ep = e.id
		LEFT JOIN cursed_technique t ON c.id = t.character_id
		WHERE c.id = ?
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("character_repo.FindAll: %w", err)
	}
	defer rows.Close()

	c := &models.Character{}
	for rows.Next() {
		var (
			ep                         int64 // first appearance episode
			e                          models.Episode
			dID, tID                   sql.NullInt64 // domain expansion & cursed technique
			dName, dDescription        sql.NullString
			tName, tDescription, tType sql.NullString
		)
		err := rows.Scan(
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
			&dID,
			&dName,
			&dDescription,
			&e.Title,
			&tID,
			&tName,
			&tDescription,
			&tType,
		)
		if err != nil {
			return nil, fmt.Errorf("character_repo.FindByID: %w", err)
		}

		if dID.Valid {
			c.Domain = &models.DomainExpansion{
				ID:          dID.Int64,
				CharacterID: c.ID,
				Name:        dName.String,
				Description: dDescription.String,
			}
		}

		if tID.Valid {
			c.CursedTechniques = append(c.CursedTechniques, models.CursedTechnique{
				ID:          tID.Int64,
				CharacterID: c.ID,
				Name:        tName.String,
				Description: tDescription.String,
				Type:        tType.String,
			})
		}

		c.FirstAppearance = &models.FirstAppearance{
			EpisodeID:    e.ID,
			EpisodeTitle: e.Title,
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("character_repo.FindByID: %w", err)
	}

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
