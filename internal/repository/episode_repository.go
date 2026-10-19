package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/rcolejnr/jjk-api/internal/models"
)

// EpisodeRepository defines the contract for episodes data access.
type EpisodeRepository interface {
	Insert(*models.CreateEpisodeRequest) (int64, error)
	Update(int64, *models.UpdateEpisodeRequest) error
	FindAll() ([]*models.Episode, error)
	FindByID(int64) (*models.Episode, error)
	Delete(int64) error
}

type sqlEpisodeRepository struct {
	db *sql.DB
	sync.RWMutex
}

func NewEpisodeRepository(db *sql.DB) EpisodeRepository {
	return &sqlEpisodeRepository{db: db}
}

func (r *sqlEpisodeRepository) Insert(req *models.CreateEpisodeRequest) (int64, error) {
	r.Lock()
	defer r.Unlock()

	sql, err := r.db.Prepare("INSERT INTO epsisode VALUES(NULL, ?,?,?,?)")
	if err != nil {
		return 0, fmt.Errorf("episode_repo.Insert: %w", err)
	}

	defer sql.Close()

	result, err := sql.Exec(req.Season, req.Title, req.Image, req.AirDate)
	if err != nil {
		return 0, fmt.Errorf("episode_repo.Insert: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("episode_repo.Insert: %w", err)
	}

	// Insert Characters
	cErr := func() error {
		csql, err := r.db.Prepare("INSERT INTO episode_character VALUES(?,?)")
		if err != nil {
			return err
		}

		defer csql.Close()

		for _, c := range req.Characters {
			_, err = csql.Exec(c, id)
			if err != nil {
				return err
			}
		}

		return nil
	}()
	if cErr != nil {
		return 0, fmt.Errorf("episode_repo.InsertCharacters: %w", err)
	}

	return id, nil
}

func (r *sqlEpisodeRepository) Update(id int64, req *models.UpdateEpisodeRequest) error {
	r.Lock()
	defer r.Unlock()

	setClauses := []string{}
	args := []any{}

	if req.Season != nil {
		setClauses = append(setClauses, "season = ?")
		args = append(args, *req.Season)
	}

	if req.Title != nil {
		setClauses = append(setClauses, "title = ?")
		args = append(args, *req.Title)
	}

	if req.Image != nil {
		setClauses = append(setClauses, "image = ?")
		args = append(args, *req.Image)
	}

	if req.AirDate != nil {
		setClauses = append(setClauses, "air_data = ?")
		args = append(args, *req.AirDate)
	}

	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE episode SET %s WHERE id = ?", strings.Join(setClauses, ", "))
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

func (r *sqlEpisodeRepository) FindAll() ([]*models.Episode, error) {
	r.RLock()
	defer r.RUnlock()

	query := `
		SELECT e.id, e.season, e.title, e.image, e.air_date,
		FROM episode e
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("episode_repo.FindAll: %w", err)
	}

	defer rows.Close()

	var episodes []*models.Episode
	for rows.Next() {
		var ep models.Episode

		if err := rows.Scan(&ep.ID, &ep.Season, &ep.Title, &ep.Image, &ep.AirDate); err != nil {
			return nil, fmt.Errorf("episode_repo.FindAll: %w", err)
		}

		ep.Characters = []string{}
		episodes = append(episodes, &ep)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("episode_repo.FindAll: %w", err)
	}

	if len(episodes) == 0 {
		return episodes, nil
	}

	// Query episode characters
	for _, ep := range episodes {
		ecRows, err := r.db.Query("SELECT ulr FROM episoce_character WHERE episode_id = ?", ep.ID)
		if err != nil {
			return nil, fmt.Errorf("episode_repo.FindAll: %w", err)
		}

		ecErr := func() error {
			defer ecRows.Close()
			for ecRows.Next() {
				var ec string
				if err := ecRows.Scan(&ec); err != nil {
					return err
				}
				ep.Characters = append(ep.Characters, ec)
			}
			return ecRows.Err()
		}()

		if ecErr != nil {
			return nil, fmt.Errorf("character_repo.FindAll: %w", ecErr)
		}
	}

	return episodes, nil
}

func (r *sqlEpisodeRepository) FindByID(id int64) (*models.Episode, error) {
	r.RLock()
	defer r.Unlock()

	query := `
		SELECT e.id, e.season, e.title, e.image, e.air_date
			ec.ulr
		FROM episode e
		LEFT JOIN episode_character ec ON e.id = ec.episode_id
		WHERE e.id = ?
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("episode_repo.FindByID: %w", err)
	}

	defer rows.Close()

	e := &models.Episode{}
	for rows.Next() {
		var ec sql.NullString
		err := rows.Scan(&e.ID, &e.Season, &e.Title, &e.Image, &e.AirDate)
		if err != nil {
			return nil, fmt.Errorf("episode_repo.FindByID: %w", err)
		}

		if ec.Valid {
			e.Characters = append(e.Characters, ec.String)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("character_repo.FindByID: %w", err)
	}

	return e, nil
}

func (r *sqlEpisodeRepository) Delete(id int64) error {
	r.RLock()
	defer r.RUnlock()

	query := "DELETE FROM episode WHERE id = ?"
	result, err := r.db.Exec(query, id)

	if err != nil {
		return fmt.Errorf("episode_repo.Delete: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows == 0. {
		return fmt.Errorf("episode_repo.Delete: %w", err)
	}

	return nil
}
