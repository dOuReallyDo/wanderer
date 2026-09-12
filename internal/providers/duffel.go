package providers

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dOuReallyDo/wanderer/internal/models"
)

// DuffelProvider — Duffel Flights API (offerte reali, prenotabili via Duffel).
// Docs: https://duffel.com/docs/api/overview
// Token: DUFFEL_ACCESS_TOKEN (duffel_test_... = sandbox simulata, duffel_live_... = offerte reali)
type DuffelProvider struct {
	token  string
	client *http.Client
}

func NewDuffelProvider() *DuffelProvider {
	return &DuffelProvider{
		token:  strings.TrimSpace(os.Getenv("DUFFEL_ACCESS_TOKEN")),
		client: &http.Client{Timeout: 25 * time.Second},
	}
}

func (d *DuffelProvider) Name() string { return "duffel" }

func (d *DuffelProvider) Available() bool {
	return strings.HasPrefix(d.token, "duffel_test_") || strings.HasPrefix(d.token, "duffel_live_")
}

// ---- strutture request/response Duffel (campi usati soltanto) ----

type duffelOfferRequest struct {
	Data duffelOfferRequestInput `json:"data"`
}

type duffelOfferRequestInput struct {
	Slices        []duffelSlice      `json:"slices"`
	Passengers    []duffelPassenger  `json:"passengers"`
	CabinClass    string             `json:"cabin_class,omitempty"`
	MaxOfferCount int                `json:"max_offer_count,omitempty"`
}

type duffelSlice struct {
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departure_date"`
}

type duffelPassenger struct {
	Type string `json:"type"`
}

type duffelAPIResponse struct {
	Data struct {
		ID     string        `json:"id"`
		Offers []duffelOffer `json:"offers"`
	} `json:"data"`
	Errors []struct {
		Title  string `json:"title"`
		Detail string `json:"detail"`
	} `json:"errors"`
}

type duffelOffer struct {
	ID            string             `json:"id"`
	TotalAmount   string             `json:"total_amount"`
	TotalCurrency string             `json:"total_currency"`
	Slices        []duffelOfferSlice `json:"slices"`
	Owner         struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
}

type duffelOfferSlice struct {
	Segments []duffelSegment `json:"segments"`
}

type duffelPlace struct {
	IataCode string `json:"iata_code"`
}

type duffelSegment struct {
	Origin      duffelPlace `json:"origin"`
	Destination duffelPlace `json:"destination"`
	DepartingAt string      `json:"departing_at"`
	ArrivingAt  string      `json:"arriving_at"`
	DurationMinutes int     `json:"duration"`
	MarketingCarrier struct {
		Name string `json:"name"`
	} `json:"marketing_carrier"`
	MarketingCarrierFlightNumber string `json:"marketing_carrier_flight_number"`
}

// Search implementa Provider.Search per Duffel.
func (d *DuffelProvider) Search(ctx context.Context, req models.SearchRequest) ([]models.Trip, error) {
	body := duffelOfferRequest{
		Data: duffelOfferRequestInput{
			Slices: []duffelSlice{{
				Origin:        req.Origin,
				Destination:   req.Destination,
				DepartureDate: req.Date,
			}},
			CabinClass:    "economy",
			MaxOfferCount: 12,
		},
	}

	nPax := req.Adults
	if nPax < 1 {
		nPax = 1
	}
	for i := 0; i < nPax; i++ {
		body.Data.Passengers = append(body.Data.Passengers, duffelPassenger{Type: "adult"})
	}

	if req.ReturnDate != "" {
		body.Data.Slices = append(body.Data.Slices, duffelSlice{
			Origin:        req.Destination,
			Destination:   req.Origin,
			DepartureDate: req.ReturnDate,
		})
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.duffel.com/air/offer_requests?return_offers=true",
		bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+d.token)
	httpReq.Header.Set("Duffel-Version", "v2")
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Accept-Encoding", "gzip")

	resp, err := d.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("duffel request: %w", err)
	}
	defer resp.Body.Close()

	// Gestione gzip esplicita (richiesta con Accept-Encoding: gzip)
	var reader io.Reader = resp.Body
	if strings.EqualFold(resp.Header.Get("Content-Encoding"), "gzip") {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("duffel gzip: %w", err)
		}
		defer gz.Close()
		reader = bufio.NewReader(gz)
	}

	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("duffel read: %w", err)
	}

	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("duffel 401: token non valido")
	}
	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return nil, fmt.Errorf("duffel %d: %s", resp.StatusCode, truncate(raw, 300))
	}

	var result duffelAPIResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("duffel decode: %w", err)
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("duffel api: %s — %s", result.Errors[0].Title, result.Errors[0].Detail)
	}

	trips := make([]models.Trip, 0, len(result.Data.Offers))
	for _, offer := range result.Data.Offers {
		if len(offer.Slices) == 0 || len(offer.Slices[0].Segments) == 0 {
			continue
		}

		var segments []models.Segment
		var totalDur int
		for _, seg := range offer.Slices[0].Segments {
			dep, _ := time.Parse(time.RFC3339, seg.DepartingAt)
			arr, _ := time.Parse(time.RFC3339, seg.ArrivingAt)
			dur := int(arr.Sub(dep).Minutes())
			if dur <= 0 {
				dur = seg.DurationMinutes
			}
			totalDur += dur
			segments = append(segments, models.Segment{
				Type:     "flight",
				FromCode: seg.Origin.IataCode,
				ToCode:   seg.Destination.IataCode,
				FromName: seg.Origin.IataCode,
				ToName:   seg.Destination.IataCode,
				Carrier:  seg.MarketingCarrier.Name,
				FlightNo: seg.MarketingCarrierFlightNumber,
				Departure: dep,
				Arrival:   arr,
				DurationMin: dur,
				Currency:  offer.TotalCurrency,
				Cabin:     "economy",
			})
		}

		price, _ := strconv.ParseFloat(offer.TotalAmount, 64)
		owner := offer.Owner.Name
		if owner == "" {
			owner = "Duffel"
		}
		trips = append(trips, models.Trip{
			ID:          "duffel-" + offer.ID,
			Segments:    segments,
			TotalPrice:  price,
			Currency:    offer.TotalCurrency,
			TotalDurMin: totalDur,
			Source:      "duffel",
			BookingURL:  bookingURLFor(req),
			Transport:   []string{"flight"},
			Notes:       "Duffel — " + owner,
		})
	}

	return trips, nil
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "…"
}

// bookingURLFor costruisce un link di prenotazione esterno per l'itinerario.
func bookingURLFor(req models.SearchRequest) string {
	q := "flights+" + req.Origin + "+to+" + req.Destination
	if req.ReturnDate != "" {
		q += "+on+" + req.Date + "+returning+" + req.ReturnDate
	} else {
		q += "+on+" + req.Date
	}
	return "https://www.google.com/travel/flights?q=" + q
}