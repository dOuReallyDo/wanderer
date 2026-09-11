package models

// City per autocomplete origine/destinazione.
type City struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"` // airport | station
	Country string `json:"country,omitempty"`
}

// SavedSearch — storico ricerche utente.
type SavedSearch struct {
	ID        string          `json:"id"`
	Request   SearchRequest   `json:"request"`
	CreatedAt string          `json:"created_at"`
	Label     string          `json:"label,omitempty"` // nickname utente
}

// FavoriteTrip — itinerari salvati dall'utente.
type FavoriteTrip struct {
	ID        string  `json:"id"`
	Trip      Trip    `json:"trip"`
	SavedAt   string  `json:"saved_at"`
	Notes     string  `json:"notes,omitempty"`
}