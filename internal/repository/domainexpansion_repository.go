package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/rcolejnr/jjk-api/models"
)

// DomainExpansionRepository defines the contrat for domain expansion data access.
type DomainExpansionRepository interface {
	Insert(*models.CreateDomainExpansionRequest) (int64, error)
	Update(int64, *models.UpdateDomainExpansionRequest) error
	FindAll() ([]*models.DomainExpansion, error)
	FindByID(int64) (*models.DomainExpansion, error)
	Delete(int64) error
}

type sqlDomainExpansionRepository struct {
	db *sql.DB
	sync.RWMutex
}

func NewDomainExpansionRepository(db *sql.DB) DomainExpansionRepository {
	return &sqlDomainExpansionRepository{db: db}
}

func (r *sqlDomainExpansionRepository) Insert(req *models.CreateDomainExpansionRequest) (int64, error) {
	r.Lock()
	defer r.Unlock()

	sql, err := r.db.Prepare("INSERT INTO domain_expansion VALUES(NULL, character_id = ?, name = ?, description = ?)")
	if err != nil {
		return 0, fmt.Errorf("domainexpansion_repository.Insert: %w", err)
	}

	defer sql.Close()

	result, err := sql.Exec(req.CharacterID, req.Name, req.Description)
	if err != nil {
		return 0, fmt.Errorf("domainexpansion_repository.Insert: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("domainexpansion_repository.Insert: %w", err)
	}

	return id, nil
}

func (r *sqlDomainExpansionRepository) Update(id int64, req *models.UpdateDomainExpansionRequest) error {
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

	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE domain_expansion SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	args = append(args, id)

	sql, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("domainexpansion_repository.Update: %w", err)
	}

	defer sql.Close()

	res, err := sql.Exec(args...)

	_, err = res.RowsAffected()

	return err
}

func (r *sqlDomainExpansionRepository) FindAll() ([]*models.DomainExpansion, error) {
	r.RLock()
	defer r.RUnlock()

	query := `
		SELECT character_id, name, description
		FROM domain_expansion
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("domainexpansion_repository.FindAll: %w", err)
	}
	defer rows.Close()

	var domains []*models.DomainExpansion
	for rows.Next() {
		var d models.DomainExpansion
		if err := rows.Scan(
			&d.ID,
			&d.CharacterID,
			&d.Name,
			&d.Description,
		); err != nil {
			return nil, fmt.Errorf("domainexpansion_repository.FindAll: %w", err)
		}
		domains = append(domains, &d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("domainexpansion_repository.FindAll: %w", err)
	}

	return domains, nil
}

func (r *sqlDomainExpansionRepository) FindByID(id int64) (*models.DomainExpansion, error) {
	r.RLock()
	defer r.RUnlock()

	query := `
		SELECT id, character_id, name, description
		FROM domain_expansion
		WHERE id = ?
	`

	var d models.DomainExpansion
	row := r.db.QueryRow(query, id)
	err := row.Scan(&d.ID, &d.CharacterID, &d.Name, &d.Description)
	if err != nil {
		return nil, fmt.Errorf("domainexpansion_repository.FindByID: %w", err)
	}

	return &d, nil
}

func (r *sqlDomainExpansionRepository) Delete(id int64) error {
	r.Lock()
	defer r.Unlock()

	query := "DELETE FROM domain_expansion WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("domainexpansion_repository.Delete: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("domainexpansion_repository.Delete: no domain found with id %d", id)
	}

	return nil
}
