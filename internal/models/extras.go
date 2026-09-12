package models

// SavedSearch — storico ricerche utente.
type SavedSearch struct {
	ID        string        `json:"id"`
	Request   SearchRequest `json:"request"`
	CreatedAt string        `json:"created_at"`
	Label     string        `json:"label,omitempty"` // nickname utente
}

// FavoriteTrip — itinerari salvati dall'utente.
type FavoriteTrip struct {
	ID      string `json:"id"`
	Trip    Trip   `json:"trip"`
	SavedAt string `json:"saved_at"`
	Notes   string `json:"notes,omitempty"`
}