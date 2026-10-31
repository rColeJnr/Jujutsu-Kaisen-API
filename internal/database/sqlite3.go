package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	DB *sql.DB
}

// NewSQLite3Repo initializes the repository, creating and configuring the database and all the necessary tables
// Database forces foreign_key.
func NewSQLite3DB(dbfile string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_foreign_keys=1&_journal_mode=WAL", dbfile))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	// CHARACTERS
	if _, err := createTableCharacter(db); err != nil {
		return nil, fmt.Errorf("Error on createTableCharacter: %w", err)
	}

	// if _, err := createTableAliases(db); err != nil {
	// 	return nil, fmt.Errorf("Error on createTableCharacter: %w", err)
	// }

	// DOMAIN EXPANSION
	if _, err := createTableDomainExpansion(db); err != nil {
		return nil, fmt.Errorf("Error on createTableDomainExpansion: %w", err)
	}

	// CURSED TECHNIQUE
	if _, err := createTableCursedTechnique(db); err != nil {
		return nil, fmt.Errorf("Error on createTableCursedTechnique: %w", err)
	}

	// EPISODES
	if _, err := createTableEpisode(db); err != nil {
		return nil, fmt.Errorf("Error on createTableEpisode: %w", err)
	}

	if _, err := createTableEpisodeCharacter(db); err != nil {
		return nil, fmt.Errorf("Error on createTableEpisodeCharacter: %w", err)
	}

	// TERMINOLOGY
	if _, err := createTableTerminology(db); err != nil {
		return nil, fmt.Errorf("Error on createTableTerminology: %w", err)
	}

	// ADMIN
	if _, err := createTableAdmin(db); err != nil {
		return nil, fmt.Errorf("Error on createTableAdmin: %w", err)
	}

	// INDEXES
	if err := createIndexes(db); err != nil {
		return nil, fmt.Errorf("Error on createIndexes: %w", err)
	}

	return &SQLiteDB{
		DB: db,
	}, nil
}

func (s *SQLiteDB) Close() error {
	return s.DB.Close()
}

// <-- TABLES -->
func createTableCharacter(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS character (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		image TEXT,
		species TEXT NOT NULL,
		status TEXT NOT NULL,
		birthday TEXT,
		age INTEGER,
		gender TEXT,
		height INTEGER,
		weight NUMERIC,
		hair_color TEXT, 
		eye_color TEXT,
		grade TEXT,
		clan TEXT,
		first_appearance_ep INTEGER NOT NULL,
		FOREIGN KEY(first_appearance_ep) REFERENCES episode(id) ON DELETE RESTRICT
	);`

	return db.Exec(sql)
}

func createTableDomainExpansion(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS domain_expansion (
		"id" INTEGER,
		"character_id" INTEGER NOT NULL UNIQUE,
		"name" TEXT NOT NULL UNIQUE,
		"description" TEXT,
		PRIMARY KEY("id"),
		FOREIGN KEY("character_id") REFERENCES character("id") ON DELETE RESTRICT
	);`

	return db.Exec(sql)
}

func createTableCursedTechnique(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS cursed_technique (
		id INTEGER,
		character_id INTEGER NOT NULL,
		name TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
		type TEXT,
		PRIMARY KEY(id),
		FOREIGN KEY(character_id) REFERENCES character(id) ON DELETE RESTRICT
	);`

	return db.Exec(sql)
}

func createTableAdmin(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS admin (
		id INTEGER,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY(id)
	);`

	return db.Exec(sql)
}

func createTableEpisode(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS episode (
		id INTEGER,
		season INTEGER NOT NULL,
		title TEXT NOT NULL UNIQUE,
		image TEXT,
		air_date TEXT,
		PRIMARY KEY(id)
	);`

	return db.Exec(sql)
}

func createTableEpisodeCharacter(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS episode_character (
		id INTEGER,
		url TEXT NOT NULL UNIQUE,
		episode_id INTEGER,
		PRIMARY KEY(id),
		FOREIGN KEY(episode_id) REFERENCES episode(id) ON DELETE CASCADE
	);`

	return db.Exec(sql)
}

func createTableTerminology(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS terminology (
		id INTEGER,
		term TEXT NOT NULL UNIQUE,
		meaning TEXT NOT NULL,
		PRIMARY KEY(id)
	)`

	return db.Exec(sql)
}

// <-- INDEXES -->
func createIndexes(db *sql.DB) error {
	queries := []string{

		`CREATE INDEX IF NOT EXISTS idx_episode_season
		ON episode (season)`, // Loading the list of episodes in a season.
		`CREATE INDEX IF NOT EXISTS idx_ep_characters
		ON episode_character (episode_id)`, // Loading the list of characters in an episode.

		`CREATE INDEX IF NOT EXISTS idx_character_species
		ON character (species)`, // Loading the list of characters where species = ?.
		`CREATE INDEX IF NOT EXISTS idx_character_status 
		ON character (status)`, // Loading the list of characters where status = ?.
		`CREATE INDEX IF NOT EXISTS idx_character_grade
		ON character (grade)`, // Loading the list of characters where grade = ?.
		`CREATE INDEX IF NOT EXISTS idx_character_first_ep
		ON character (first_appearance_ep)`, // Enables fast checks on character table on delete.

		`CREATE INDEX IF NOT EXISTS idx_character_ctech
		ON cursed_technique (character_id)`, // Loading a character cursed technique.
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("Failed to create index: %w", err)
		}
	}

	return nil
}
