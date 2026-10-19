package models

type Character struct {
	ID               int64             `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Image            string            `json:"images,omitempty"`
	Species          string            `json:"species"`
	Status           string            `json:"status"`
	Birthday         string            `json:"birthday"`
	Age              int               `json:"age"`
	Gender           string            `json:"gender"`
	Height           int               `json:"height"`
	Weight           float64           `json:"weight"`
	HairColor        string            `json:"hair_color"`
	EyeColor         string            `json:"eye_color"`
	Grade            string            `json:"grade"`
	Clan             string            `json:"clan"`
	CursedTechniques []CursedTechnique `json:"cursed_techniques,omitempty"`
	Domain           *DomainExpansion  `json:"domain,omitempty"`
	FirstAppearance  *FirstAppearance  `json:"first_appearance,omitempty"`
}

type FirstAppearance struct {
	EpisodeID    int64
	EpisodeTitle string
}

type CreateCharacterRequest struct {
	Name            string
	Description     string
	Image           string
	Species         string
	Status          string
	Birthday        string
	Age             int
	Gender          string
	Height          int
	Weight          float64
	HairColor       string
	EyeColor        string
	Grade           string
	Clan            string
	FirstAppearance int64
}

type UpdateCharacterRequest struct {
	Name        *string
	Description *string
	Image       *string
	Species     *string
	Status      *string
	Birthday    *string
	Age         *int
	Gender      *string
	Height      *int
	Weight      *float64
	HairColor   *string
	EyeColor    *string
	Grade       *string
	Clan        *string
}

// CharacterPaginatedResponse wraps a list of characters with pagination metadata.
type CharacterPaginatedResponse struct {
	Data  []Character `json:"data"`
	Count int         `json:"count"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}
