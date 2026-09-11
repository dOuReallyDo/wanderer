package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/dOuReallyDo/wanderer/internal/models"
)

// KiwiProvider — Tequila by Kiwi.com (low-cost, offerte reali).
// Docs: https://tequila.kiwi.com/portal/docs/tequila/
type KiwiProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewKiwiProvider() *KiwiProvider {
	return &KiwiProvider{
		apiKey:  os.Getenv("TEQUILA_API_KEY"),
		baseURL: envOr("TEQUILA_BASE_URL", "https://tequila-api.kiwi.com"),
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (k *KiwiProvider) Name() string { return "kiwi-tequila" }
func (k *KiwiProvider) Available() bool {
	return k.apiKey != ""
}

func (k *KiwiProvider) Search(ctx context.Context, req models.SearchRequest) ([]models.Trip, error) {
	params := url.Values{}
	params.Set("fly_from", req.Origin)
	params.Set("fly_to", req.Destination)
	params.Set("date_from", req.Date)
	params.Set("date_to", req.Date)
	params.Set("adults", strconv.Itoa(max(req.Adults, 1)))
	params.Set("curr", envOrDefault(req.Currency, "EUR"))
	params.Set("limit", "10")
	if req.MaxPrice > 0 {
		params.Set("price_to", fmt.Sprintf("%.0f", req.MaxPrice))
	}
	if req.ReturnDate != "" {
		params.Set("return_from", req.ReturnDate)
		params.Set("return_to", req.ReturnDate)
	}

	apiURL := k.baseURL + "/v2/search?" + params.Encode()
	httpReq, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	httpReq.Header.Set("apikey", k.apiKey)

	resp, err := k.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("kiwi search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("kiwi %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []struct {
			ID           string      `json:"id"`
			Price        float64     `json:"price"`
			Currency     string      `json:"currency"`
			Duration     struct {
				Total int `json:"total"`
			} `json:"duration"`
			Route []struct {
				From    string `json:"flyFrom"`
				To      string `json:"flyTo"`
				CityFrom string `json:"cityFrom"`
				CityTo   string `json:"cityTo"`
				Airline string `json:"airline"`
				FlightNo int    `json:"flight_no"`
				UTCDepart  int64 `json:"dTime"`
				UTCArrival int64 `json:"aTime"`
			} `json:"route"`
			DeepLink string `json:"deep_link"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("kiwi decode: %w", err)
	}

	var trips []models.Trip
	for _, offer := range result.Data {
		var segments []models.Segment
		for _, seg := range offer.Route {
			dep := time.Unix(seg.UTCDepart, 0)
			arr := time.Unix(seg.UTCArrival, 0)
			segments = append(segments, models.Segment{
				Type:       "flight",
				FromCode:   seg.From,
				ToCode:     seg.To,
				FromName:   seg.CityFrom,
				ToName:     seg.CityTo,
				Carrier:    seg.Airline,
				FlightNo:   seg.Airline + "-" + strconv.Itoa(seg.FlightNo),
				Departure:  dep,
				Arrival:    arr,
				DurationMin: int(arr.Sub(dep).Minutes()),
				BookingURL: offer.DeepLink,
			})
		}
		trips = append(trips, models.Trip{
			ID:          "kiwi-" + offer.ID,
			Segments:    segments,
			TotalPrice:  offer.Price,
			Currency:    offer.Currency,
			TotalDurMin: offer.Duration.Total / 60,
			Source:      "kiwi-tequila",
			BookingURL:  offer.DeepLink,
			Transport:   []string{"flight"},
		})
	}
	return trips, nil
}