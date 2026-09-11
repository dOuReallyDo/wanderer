package models

import "time"

// Trip rappresenta un itinerario completo proposto all'utente.
type Trip struct {
	ID          string    `json:"id"`
	Segments    []Segment `json:"segments"`
	TotalPrice  float64   `json:"total_price"`
	Currency    string    `json:"currency"`
	TotalDurMin int       `json:"total_duration_min"`
	Score       float64   `json:"score"`        // Wander Score 0-100
	Rank        int       `json:"rank"`         // posizione in classifica
	Source      string    `json:"source"`       // provider che ha generato l'offerta
	BookingURL  string    `json:"booking_url"`  // link reale per sottoscrivere
	Transport   []string  `json:"transport_types"` // "flight", "train", "mixed"
	Notes       string    `json:"notes,omitempty"`
}

// Segment è una tratta elementare dell'itinerario.
type Segment struct {
	Type        string    `json:"type"` // flight | train
	FromCode    string    `json:"from_code"`
	ToCode      string    `json:"to_code"`
	FromName    string    `json:"from_name"`
	ToName      string    `json:"to_name"`
	Carrier     string    `json:"carrier"`
	FlightNo    string    `json:"flight_no,omitempty"`
	TrainNo     string    `json:"train_no,omitempty"`
	Departure   time.Time `json:"departure"`
	Arrival     time.Time `json:"arrival"`
	DurationMin int       `json:"duration_min"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	BookingURL  string    `json:"booking_url"`
	Cabin       string    `json:"cabin,omitempty"`
}

// SearchRequest è la query dell'utente.
type SearchRequest struct {
	Origin          string  `json:"origin"`           // IATA o stazione
	Destination     string  `json:"destination"`
	Date            string  `json:"date"`             // YYYY-MM-DD
	ReturnDate      string  `json:"return_date,omitempty"`
	Adults          int     `json:"adults"`
	MaxPrice        float64 `json:"max_price"`
	Currency        string  `json:"currency"`        // EUR default
	MaxDurationH    int     `json:"max_duration_h"`
	EarliestDepart  string  `json:"earliest_depart,omitempty"` // HH:MM
	LatestDepart    string  `json:"latest_depart,omitempty"`
	PreferredTransport string `json:"preferred_transport,omitempty"` // flight | train | any
}

// SearchResponse è il risultato con ranking.
type SearchResponse struct {
	Request   SearchRequest `json:"request"`
	Trips     []Trip        `json:"trips"`
	Best      *Trip         `json:"best,omitempty"`
	Count     int           `json:"count"`
	SearchedAt time.Time   `json:"searched_at"`
	TookMs    int64         `json:"took_ms"`
	Providers []string     `json:"providers_queried"`
}

// RankWeights definisce i pesi del Wander Score.
type RankWeights struct {
	Cost      float64 `json:"cost"`       // peso costo (default 40)
	Duration  float64 `json:"duration"`   // peso durata (default 25)
	Schedule  float64 `json:"schedule"`   // peso orario (default 20)
	Serendipity float64 `json:"serendipity"` // peso sorpresa/offerta (default 15)
}

func DefaultWeights() RankWeights {
	return RankWeights{
		Cost:       40,
		Duration:   25,
		Schedule:   20,
		Serendipity: 15,
	}
}