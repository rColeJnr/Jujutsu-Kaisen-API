package models

type Location struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	LocatedIn  string   `json:"located_in"`
	Episodes   []string `json:"episodes,omitempty"`
	Characters []string `json:"characters,omitempty"`
}
