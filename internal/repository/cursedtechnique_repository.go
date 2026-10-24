package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/rcolejnr/jjk-api/internal/models"
)

// CursedTechniqueRepository defines the contrat for CursedTechnique data access.
type CursedTechniqueRepository interface {
	Insert(*models.CreateCursedTechniqueRequest) (int64, error)
	Update(int64, *models.UpdateCursedTechniqueRequest) error
	FindAll(*models.CursedTechniqueFilter) ([]*models.CursedTechnique, error)
	FindByID(int64) (*models.CursedTechnique, error)
	Delete(int64) error
}

type sqlCursedTechniqueRepository struct {
	db *sql.DB
	sync.RWMutex
}

func NewCursedTechniqueRepository(db *sql.DB) CursedTechniqueRepository {
	return &sqlCursedTechniqueRepository{db: db}
}

func (r *sqlCursedTechniqueRepository) Insert(req *models.CreateCursedTechniqueRequest) (int64, error) {
	r.Lock()
	defer r.Unlock()

	sql, err := r.db.Prepare("INSERT INTO cursed_technique (character_id, name, description, type) VALUES(?,?,?,?)")
	if err != nil {
		return 0, fmt.Errorf("cursedtechnique_repository.Insert: here %w", err)
	}

	defer sql.Close()

	result, err := sql.Exec(req.CharacterID, req.Name, req.Description, req.Type)
	if err != nil {
		return 0, fmt.Errorf("cursedtechnique_repository.Insert: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("cursedtechnique_repository.Insert: %w", err)
	}

	return id, nil
}

func (r *sqlCursedTechniqueRepository) Update(id int64, req *models.UpdateCursedTechniqueRequest) error {
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

	if req.Type != nil {
		setClauses = append(setClauses, "type = ?")
		args = append(args, *req.Type)
	}

	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE cursed_technique SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	args = append(args, id)

	sql, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("cursedtechnique_repository.Update: %w", err)
	}

	defer sql.Close()

	res, err := sql.Exec(args...)

	_, err = res.RowsAffected()

	return err
}

func (r *sqlCursedTechniqueRepository) FindAll(f *models.CursedTechniqueFilter) ([]*models.CursedTechnique, error) {
	r.RLock()
	defer r.RUnlock()

	query := `
		SELECT id, character_id, name, description, type
		FROM cursed_technique
	`

	if f.Name != "" {
		query += ` WHERE name LIKE ? ESCAPE '\'`
	}

	rows, err := r.db.Query(query, "%"+escapeLike(f.Name)+"%")
	if err != nil {
		return nil, fmt.Errorf("cursedtechnique_repository.FindAll: %w", err)
	}
	defer rows.Close()

	var cursedTechniques []*models.CursedTechnique
	for rows.Next() {
		var ct models.CursedTechnique
		if err := rows.Scan(
			&ct.ID,
			&ct.CharacterID,
			&ct.Name,
			&ct.Description,
			&ct.Type,
		); err != nil {
			return nil, fmt.Errorf("cursedtechnique_repository.FindAll: %w", err)
		}
		cursedTechniques = append(cursedTechniques, &ct)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cursedtechnique_repository.FindAll: %w", err)
	}

	return cursedTechniques, nil
}

func (r *sqlCursedTechniqueRepository) FindByID(id int64) (*models.CursedTechnique, error) {
	r.RLock()
	defer r.RUnlock()

	query := `
		SELECT id, character_id, name, description, type
		FROM cursed_technique
		WHERE id = ?
	`

	var d models.CursedTechnique
	row := r.db.QueryRow(query, id)
	err := row.Scan(&d.ID, &d.CharacterID, &d.Name, &d.Description, &d.Type)
	if err != nil {
		return nil, fmt.Errorf("cursedtechnique_repository.FindByID: %w", err)
	}

	return &d, nil
}

func (r *sqlCursedTechniqueRepository) Delete(id int64) error {
	r.Lock()
	defer r.Unlock()

	query := "DELETE FROM cursed_technique WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("cursedtechnique_repository.Delete: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("cursedtechnique_repository.Delete: no cursed technique found with id %d", id)
	}

	return nil
}
