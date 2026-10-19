package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rcolejnr/jjk-api/models"
)

type sqLiteDB struct {
	db *sql.DB
}

// NewSQLite3Repo initializes the repository, creating and configuring the database and all the necessary tables
// Database forces foreign_key.
func NewSQLite3DB(dbfile string) (*sqLiteDB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_pragma=foreign_keys(1)", dbfile))
	if err != nil {
		return nil, err
	}
	defer db.Close() // Need to remove this line, to keep db open during application runtime.

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	// CHARACTERS
	if _, err := CreateTableCharacter(db); err != nil {
		return nil, err
	}

	if _, err := CreateTableDomainExpansion(db); err != nil {
		return nil, err
	}

	if _, err := CreateTableCursedTechnique(db); err != nil {
		return nil, err
	}

	if _, err := CreateTableAliases(db); err != nil {
		return nil, err
	}

	// EPISODES
	if _, err := models.CreateTableEpisode(db); err != nil {
		return nil, err
	}

	if _, err := models.CreateTableEpisodeCharacter(db); err != nil {
		return nil, err
	}

	return &sqLiteDB{
		db: db,
	}, nil
}

func (s *sqLiteDB) Close() error {
	return s.db.Close()
}

func CreateTableCharacter(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS character (
		"id" INTEGER PRIMARY KEY,
		"name" TEXT NOT NULL UNIQUE,
		"description" TEXT NOT NULL,
		"image" TEXT,
		"species" TEXT NOT NULL,
		"status" TEXT NOT NULL,
		"birthday" TEXT NOT NULL,
		"age" INTERGER NOT NULL,
		"gender" TEXT NOT NULL,
		"height" INTEGER NOT NULL,
		"weight" NUMERIC NOT NUL,
		"hair_color" TEXT NOT NULL, 
		"eye_color" TEXT NOT NULL,
		"grade" TEXT NOT NULL,
		"clan" TEXT,
		"first_appearance_ep" INTEGER NOT NULL,
		"FOREIGN KEY("first_appearance_ep") REFERENCES "episode"("id") ON DELETE RESTRICT 
	);`

	return db.Exec(sql)
}

func CreateTableDomainExpansion(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS domain_expansion (
		"id" INTEGER PRIMARY KEY,
		"character_id" INTEGER NOT NULL,
		"name" TEXT NOT NULL,
		"description" TEXT,
		PRIMARY KEY("character_id", "id"),
		FOREIGN KEY("character_id") REFERENCES character("id") ON DELETE CASCADE
	);`

	return db.Exec(sql)
}

func CreateTableCursedTechnique(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS cursed_technique (
		id INTEGER PRIMARY KEY,
		character_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		type TEXT,
		PRIMARY KEY(character_id, id),
		FOREIGN KEY(character_id) REFERENCES character(id) ON DELETE CASCADE
	);`

	return db.Exec(sql)
}

func CreateTableAliases(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS aliases (
		alias TEXT NOT NULL,
		character_id INTEGER,
		PRIMARY KEY(character_id, alias),
		FOREIGN KEY(character_id) REFERENCES character(id) ON DELETE CASCADE
	);`

	return db.Exec(sql)
}
