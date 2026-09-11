package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dOuReallyDo/wanderer/internal/models"
)

// AmadeusProvider — Amadeus Self-Service API (voli reali, prenotabili).
// Docs: https://developers.amadeus.com/self-service/category/air
type AmadeusProvider struct {
	apiKey    string
	apiSecret string
	baseURL   string
	token     string
	tokenExp  time.Time
	client    *http.Client
}

func NewAmadeusProvider() *AmadeusProvider {
	return &AmadeusProvider{
		apiKey:    os.Getenv("AMADEUS_API_KEY"),
		apiSecret: os.Getenv("AMADEUS_API_SECRET"),
		baseURL:   envOr("AMADEUS_BASE_URL", "https://test.api.amadeus.com"),
		client:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (a *AmadeusProvider) Name() string { return "amadeus" }
func (a *AmadeusProvider) Available() bool {
	return a.apiKey != "" && a.apiSecret != ""
}

func (a *AmadeusProvider) refreshToken(ctx context.Context) error {
	if a.token != "" && time.Now().Before(a.tokenExp) {
		return nil
	}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", a.apiKey)
	data.Set("client_secret", a.apiSecret)

	req, _ := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/v1/security/oauth/token",
		strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("amadeus auth: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("amadeus auth %d: %s", resp.StatusCode, string(body))
	}

	var tr struct {
		AccessToken string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return fmt.Errorf("amadeus auth decode: %w", err)
	}
	a.token = tr.AccessToken
	a.tokenExp = time.Now().Add(time.Duration(tr.ExpiresIn-60) * time.Second)
	return nil
}

func (a *AmadeusProvider) Search(ctx context.Context, req models.SearchRequest) ([]models.Trip, error) {
	if err := a.refreshToken(ctx); err != nil {
		return nil, err
	}

	// Amadeus Flight Offers Search v2
	params := url.Values{}
	params.Set("originLocationCode", req.Origin)
	params.Set("destinationLocationCode", req.Destination)
	params.Set("departureDate", req.Date)
	params.Set("adults", strconv.Itoa(max(req.Adults, 1)))
	params.Set("currencyCode", envOrDefault(req.Currency, "EUR"))
	params.Set("max", "10")
	if req.MaxPrice > 0 {
		params.Set("maxPrice", fmt.Sprintf("%.0f", req.MaxPrice))
	}
	if req.ReturnDate != "" {
		params.Set("returnDate", req.ReturnDate)
	}

	apiURL := a.baseURL + "/v2/shopping/flight-offers?" + params.Encode()
	httpReq, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	httpReq.Header.Set("Authorization", "Bearer "+a.token)

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("amadeus search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("amadeus search %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []struct {
			ID           string `json:"id"`
			PriceMaps    struct {
				Currency string `json:"currency"`
				GrandTotal string `json:"grandTotal"`
			} `json:"price"`
			Itineraries []struct {
				Segments []struct {
					Departure struct {
						IataCode string `json:"iataCode"`
						At       string `json:"at"`
					} `json:"departure"`
					Arrival struct {
						IataCode string `json:"iataCode"`
						At       string `json:"at"`
					} `json:"arrival"`
					CarrierCode string `json:"carrierCode"`
					Number      string `json:"number"`
					Duration    string `json:"duration"`
					Aircraft    struct{ Code string `json:"code"` } `json:"aircraft"`
				} `json:"segments"`
			} `json:"itineraries"`
		} `json:"data"`
		Dictionaries struct {
			Carriers map[string]string `json:"carriers"`
		} `json:"dictionaries"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&result); err != nil {
		return nil, fmt.Errorf("amadeus decode: %w", err)
	}

	var trips []models.Trip
	for _, offer := range result.Data {
		var segments []models.Segment
		var totalDur int
		var earliestDep, latestArr time.Time

		for _, itin := range offer.Itineraries {
			for _, seg := range itin.Segments {
				dep, _ := time.Parse(time.RFC3339, seg.Departure.At)
				arr, _ := time.Parse(time.RFC3339, seg.Arrival.At)
				dur := int(arr.Sub(dep).Minutes())
				totalDur += dur
				if earliestDep.IsZero() || dep.Before(earliestDep) {
					earliestDep = dep
				}
				if arr.After(latestArr) {
					latestArr = arr
				}
				segments = append(segments, models.Segment{
					Type:        "flight",
					FromCode:    seg.Departure.IataCode,
					ToCode:      seg.Arrival.IataCode,
					FromName:    seg.Departure.IataCode,
					ToName:      seg.Arrival.IataCode,
					Carrier:     result.Dictionaries.Carriers[seg.CarrierCode],
					FlightNo:    seg.CarrierCode + seg.Number,
					Departure:   dep,
					Arrival:     arr,
					DurationMin: dur,
					BookingURL:  "https://www.amadeus.com",
					Cabin:       "ECONOMY",
				})
			}
		}

		price, _ := strconv.ParseFloat(offer.PriceMaps.GrandTotal, 64)

		trips = append(trips, models.Trip{
			ID:          "amadeus-" + offer.ID,
			Segments:    segments,
			TotalPrice:  price,
			Currency:    offer.PriceMaps.Currency,
			TotalDurMin: totalDur,
			Source:      "amadeus",
			BookingURL:  "https://www.amadeus.com",
			Transport:   []string{"flight"},
		})
	}

	return trips, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
func envOrDefault(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}