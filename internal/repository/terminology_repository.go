package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/rcolejnr/jjk-api/internal/models"
)

type TerminologyRepository interface {
	Insert(*models.CreateTerminologyRequest) (int64, error)
	Update(int64, *models.UpdateTerminologyRequest) error
	FindAll(*models.TerminologyFilter) ([]*models.Terminology, error)
	FindByID(int64) (*models.Terminology, error)
}

type sqlTerminologyRepository struct {
	db *sql.DB
	sync.RWMutex
}

func NewTerminologyRepository(db *sql.DB) TerminologyRepository {
	return &sqlTerminologyRepository{db: db}
}

func (r *sqlTerminologyRepository) Insert(req *models.CreateTerminologyRequest) (int64, error) {
	r.Lock()
	defer r.Unlock()

	sql, err := r.db.Prepare("INSERT INTO terminology VALUES(NULL, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("terminology_repository.Insert: %w", err)
	}

	defer sql.Close()

	result, err := sql.Exec(req.Term, req.Meaning)
	if err != nil {
		return 0, fmt.Errorf("terminology_repository.Insert: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("terminology_repository.Insert: %w", err)
	}

	return id, nil
}
func (r *sqlTerminologyRepository) Update(id int64, req *models.UpdateTerminologyRequest) error {
	r.Lock()
	defer r.Unlock()

	var setClauses []string
	var args []any

	if req.Term != nil {
		setClauses = append(setClauses, "term = ?")
		args = append(args, *req.Term)
	}

	if req.Meaning != nil {
		setClauses = append(setClauses, "meaning = ?")
		args = append(args, *req.Meaning)
	}

	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, id)

	query := fmt.Sprintf("UPDATE terminology SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("terminology_repository.Update: %w", err)
	}

	result, err := stmt.Exec(args...)
	if err != nil {
		return fmt.Errorf("terminology_repository.Update: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("terminology_repository.Update: no terminologies found with id %d", id)
	}

	return nil
}

func (r *sqlTerminologyRepository) FindAll(f *models.TerminologyFilter) ([]*models.Terminology, error) {
	r.RLock()
	defer r.RUnlock()

	query := "SELECT id, term, meaning FROM terminology"

	if f.Term != "" {
		query += ` WHERE term LIKE ? ESCAPE '\'`
	}

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("terminology_repository.FindAll: %w", err)
	}

	defer stmt.Close()

	rows, err := r.db.Query(query, "%"+escapeLike(f.Term)+"%")
	if err != nil {
		return nil, fmt.Errorf("terminology_repository.FindAll: %w", err)
	}

	var terms []*models.Terminology
	for rows.Next() {
		var t models.Terminology
		if err = rows.Scan(&t.ID, &t.Term, &t.Meaning); err != nil {
			return nil, fmt.Errorf("terminology_repository.FindAll: %w", err)
		}

		terms = append(terms, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("terminology_repository.FindAll: %w", err)
	}

	return terms, nil
}
func (r *sqlTerminologyRepository) FindByID(id int64) (*models.Terminology, error) {
	r.RLock()
	defer r.RUnlock()

	query := "SELECT id, term, meaning FROM terminology WHERE id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("terminology_repository.FindByID: %w", err)
	}

	defer stmt.Close()

	var t models.Terminology
	row := r.db.QueryRow(query, id)
	if err = row.Scan(&t.ID, &t.Term, &t.Meaning); err != nil {
		return nil, fmt.Errorf("terminology_repository.FindByID: %w", err)
	}

	return &t, nil
}
