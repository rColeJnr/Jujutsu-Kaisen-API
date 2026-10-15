package models

type Episode struct {
	ID         int64    `json:"id"`
	Season     int      `json:"season"`
	Title      string   `json:"title"`
	Image      string   `json:"image"`
	AirDate    string   `json:"air_date"`
	Locations  []string `json:"locations,omitempty"`
	Characters []string `json:"characters,omitempty"`
}
