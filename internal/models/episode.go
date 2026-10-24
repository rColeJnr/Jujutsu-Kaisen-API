package models

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

type EpisodeFilter struct {
	Season int
	Title  string
}
