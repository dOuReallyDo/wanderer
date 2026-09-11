package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dOuReallyDo/wanderer/internal/models"
)

// DBProvider — Deutsche Bahn / Trainline-style API per treni europei.
// Usa l'API pubblica DB Fahrplan (non richiede chiave per uso base).
// Docs: https://developers.deutschebahn.com/
type DBProvider struct {
	baseURL string
	client  *http.Client
}

func NewDBProvider() *DBProvider {
	return &DBProvider{
		baseURL: envOr("DB_BASE_URL", "https://v6.db.transport.rest"),
		client: &http.Client{Timeout: 5 * time.Second}, // short: API pubblica può essere down
	}
}

func (d *DBProvider) Name() string { return "deutsche-bahn" }
func (d *DBProvider) Available() bool { return true } // API pubblica, no key

func (d *DBProvider) Search(ctx context.Context, req models.SearchRequest) ([]models.Trip, error) {
	// DB Fahrplan API: /trips?from=...&to=...&departure=...
	params := url.Values{}
	params.Set("from", req.Origin)
	params.Set("to", req.Destination)
	params.Set("departure", req.Date+"T00:00")
	params.Set("results", "10")

	apiURL := d.baseURL + "/trips?" + params.Encode()
	httpReq, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)

	resp, err := d.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("db search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("db %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Trips []struct {
			ID       string `json:"id"`
			Origin   struct {
				Name      string    `json:"name"`
				Departure time.Time `json:"departure"`
			} `json:"origin"`
			Destination struct {
				Name     string    `json:"name"`
				Arrival  time.Time `json:"arrival"`
			} `json:"destination"`
			Duration int `json:"duration"` // secondi
			Price    struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
			} `json:"price"`
			Legs []struct {
				Origin      string    `json:"origin"`
				Destination string    `json:"destination"`
				Departure   time.Time `json:"departure"`
				Arrival     time.Time `json:"arrival"`
				LineName    string    `json:"lineName"`
				Direction   string    `json:"direction"`
			} `json:"legs"`
		} `json:"trips"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("db decode: %w", err)
	}

	var trips []models.Trip
	for _, t := range result.Trips {
		var segments []models.Segment
		for _, leg := range t.Legs {
			dur := int(leg.Arrival.Sub(leg.Departure).Minutes())
			segments = append(segments, models.Segment{
				Type:       "train",
				FromName:   leg.Origin,
				ToName:     leg.Destination,
				Carrier:    "Deutsche Bahn",
				TrainNo:    leg.LineName,
				Departure:  leg.Departure,
				Arrival:    leg.Arrival,
				DurationMin: dur,
				BookingURL:  "https://www.bahn.com",
			})
		}
		trips = append(trips, models.Trip{
			ID:          "db-" + t.ID,
			Segments:    segments,
			TotalPrice:  t.Price.Amount,
			Currency:    t.Price.Currency,
			TotalDurMin: t.Duration / 60,
			Source:      "deutsche-bahn",
			BookingURL:  "https://www.bahn.com",
			Transport:   []string{"train"},
		})
	}
	return trips, nil
}

// suppress unused
var _ = strconv.Itoa