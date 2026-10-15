package models

type Terminology struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Episodes    []string `json:"episodes,omitempty"`
	Characters  []string `json:"characters,omitempty"`
}
