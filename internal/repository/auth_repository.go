package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rcolejnr/jjk-api/internal/models"
)

type AuthRepository interface {
	Insert(string, string) (*models.Admin, error)
	Update(string, string) error
	FindByEmail(string) (*models.Admin, error)
	FindByID(int64) (*models.Admin, error)
}

type sqlAuthRepository struct {
	db *sql.DB
	sync.RWMutex
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &sqlAuthRepository{db: db}
}

func (r *sqlAuthRepository) Insert(email, password string) (*models.Admin, error) {
	r.Lock()
	defer r.Unlock()

	query := `
		INSERT INTO admin (email, password, created_at)
		VALUES(?,?,?)
		RETURNING id, email, created_at
	`

	var admin models.Admin
	if err := r.db.QueryRow(query, email, password, time.Now()).Scan(
		&admin.ID, &admin.Email, &admin.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("authrepository.Insert: %w", err)
	}

	return &admin, nil
}

func (r *sqlAuthRepository) Update(email, password string) error {
	r.Lock()
	defer r.Unlock()

	query := `UPDATE admin SET password = ? WHERE email = ?`

	res, err := r.db.Exec(query, password, email)
	if err != nil {
		return fmt.Errorf("authRepository.Update: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("authRepository.Update: admin not found.")
	}

	return nil
}

func (r *sqlAuthRepository) FindByEmail(email string) (*models.Admin, error) {
	r.RLock()
	defer r.RUnlock()

	query := `SELECT id, email, password FROM admin WHERE email = ?`

	var admin models.Admin
	if err := r.db.QueryRow(query, email).Scan(&admin.ID, &admin.Email, &admin.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("authrepository.FindByEmail: %w", err)
	}

	return &admin, nil
}

func (r *sqlAuthRepository) FindByID(id int64) (*models.Admin, error) {
	r.RLock()
	defer r.RUnlock()

	query := `SELECT id, email, created_at FROM admin WHERE id = ?`

	var admin models.Admin
	if err := r.db.QueryRow(query, id).Scan(
		&admin.ID, &admin.Email, &admin.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("authrepository.FindByID: %w", err)
	}

	return &admin, nil
}
