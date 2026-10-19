package models

import "database/sql"

type Episode struct {
	ID         int64    `json:"id"`
	Season     int64    `json:"season"`
	Title      string   `json:"title"`
	Image      string   `json:"image"`
	AirDate    string   `json:"air_date"`
	Characters []string `json:"characters,omitempty"`
}

type EpisodeCharacter struct {
	Url       string
	EpisodeID int64
}

type CreateEpisodeRequest struct {
	Season     int64    `json:"season"`
	Title      string   `json:"title"`
	Image      string   `json:"image"`
	AirDate    string   `json:"air_date"`
	Characters []string `json:"characters,omitempty"`
}

type UpdateEpisodeRequest struct {
	Season  *int64  `json:"season"`
	Title   *string `json:"title"`
	Image   *string `json:"image"`
	AirDate *string `json:"air_date"`
}

func CreateTableEpisode(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS episode (
		id INTEGER PRIMARY KEY,
		season INTEGER NOT NULL,
		title TEXT NOT NULL,
		image TEXT,
		air_date TEXT,
	);`

	return db.Exec(sql)
}

func CreateTableEpisodeCharacter(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS episode_character (
		url TEXT NOT NULL,
		episode_id INTEGER,
		PRIMARY KEY(episode_id, url),
		FOREIGN KEY(episode_id) REFERENCES episode(id) ON DELETE CASCADE
	);`

	return db.Exec(sql)
}
