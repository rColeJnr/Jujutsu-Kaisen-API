package models

type Character struct {
	ID               int64              `json:"id"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	Image            string             `json:"images,omitempty"`
	Species          string             `json:"species"`
	Status           string             `json:"status"`
	Birthday         string             `json:"birthday"`
	Age              int                `json:"age"`
	Gender           string             `json:"gender"`
	Height           int                `json:"height"`
	Weight           float64            `json:"weight"`
	HairColor        string             `json:"hair_color"`
	EyeColor         string             `json:"eye_color"`
	Grade            string             `json:"grade"`
	Clan             string             `json:"clan"`
	CursedTechniques []*CursedTechnique `json:"cursed_techniques,omitempty"`
	Domain           *DomainExpansion   `json:"domain_expansion,omitempty"`
	FirstAppearance  *FirstAppearance   `json:"first_appearance,omitempty"`
}

type FirstAppearance struct {
	EpisodeID    int64  `json:"episode_id"`
	EpisodeTitle string `json:"episode_title"`
}

type CreateCharacterRequest struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Image             string  `json:"images,omitempty"`
	Species           string  `json:"species"`
	Status            string  `json:"status"`
	Birthday          string  `json:"birthday"`
	Age               int     `json:"age"`
	Gender            string  `json:"gender"`
	Height            int     `json:"height"`
	Weight            float64 `json:"weight"`
	HairColor         string  `json:"hair_color"`
	EyeColor          string  `json:"eye_color"`
	Grade             string  `json:"grade"`
	Clan              string  `json:"clan"`
	FirstAppearanceEp int64   `json:"first_appearance_ep,omitempty"`
}

type UpdateCharacterRequest struct {
	Name              *string  `json:"name"`
	Description       *string  `json:"description"`
	Image             *string  `json:"images,omitempty"`
	Species           *string  `json:"species"`
	Status            *string  `json:"status"`
	Birthday          *string  `json:"birthday"`
	Age               *int     `json:"age"`
	Gender            *string  `json:"gender"`
	Height            *int     `json:"height"`
	Weight            *float64 `json:"weight"`
	HairColor         *string  `json:"hair_color"`
	EyeColor          *string  `json:"eye_color"`
	Grade             *string  `json:"grade"`
	Clan              *string  `json:"clan"`
	FirstAppearanceEp *int64   `json:"first_appearance_ep,omitempty"`
}

// CharacterPaginatedResponse wraps a list of characters with pagination metadata.
type CharacterPaginatedResponse struct {
	Data  []Character `json:"data"`
	Count int         `json:"count"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

// CharacterFilter defines the filters to http.MethodGet characters
type CharacterFilter struct {
	Name, Species, Status, Gender, Grade string
}
