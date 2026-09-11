package providers

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/dOuReallyDo/wanderer/internal/models"
)

// DemoProvider genera offerte realistiche per demo/testing.
// Simula voli e treni con prezzi e orari plausibili.
// Disattivato automaticamente quando provider reali sono attivi.
type DemoProvider struct {
	seed *rand.Rand
}

func NewDemoProvider() *DemoProvider {
	return &DemoProvider{seed: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (d *DemoProvider) Name() string { return "demo" }
func (d *DemoProvider) Available() bool { return true }

func (d *DemoProvider) Search(ctx context.Context, req models.SearchRequest) ([]models.Trip, error) {
	// Genera 8-12 offerte miste volo/treno
	numOffers := 8 + d.seed.Intn(5)
	trips := make([]models.Trip, 0, numOffers)

	// Parse date
	baseDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("demo: invalid date: %w", err)
	}

	carriers := []string{"Ryanair", "EasyJet", "Lufthansa", "ITA Airways", "Vueling", "KLM"}
	trainCarriers := []string{"Trenitalia", "Deutsche Bahn", "SNCF", "Trenitalia Frecciarossa", "ÖBB"}

	for i := 0; i < numOffers; i++ {
		isTrain := d.seed.Intn(3) == 0 // ~33% treni
		isMixed := d.seed.Intn(5) == 0 // ~20% misto

		// Orario partenza: 6:00 - 21:00
		depHour := 6 + d.seed.Intn(15)
		depMin := d.seed.Intn(12) * 5
		depTime := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(),
			depHour, depMin, 0, 0, time.Local)

		// Durata: volo 1-5h, treno 3-12h
		var durMin int
		if isTrain {
			durMin = 180 + d.seed.Intn(540)
		} else {
			durMin = 60 + d.seed.Intn(240)
		}
		arrTime := depTime.Add(time.Duration(durMin) * time.Minute)

		// Prezzo: volo 30-400, treno 15-200
		var price float64
		if isTrain {
			price = 15 + d.seed.Float64()*185
		} else {
			price = 30 + d.seed.Float64()*370
		}

		carrier := ""
		transportType := "flight"
		if isTrain {
			carrier = trainCarriers[d.seed.Intn(len(trainCarriers))]
			transportType = "train"
		} else {
			carrier = carriers[d.seed.Intn(len(carriers))]
		}
		if isMixed {
			transportType = "mixed"
		}

		seg := models.Segment{
			Type:       transportType,
			FromCode:   req.Origin,
			ToCode:     req.Destination,
			FromName:   req.Origin,
			ToName:     req.Destination,
			Carrier:    carrier,
			FlightNo:   fmt.Sprintf("%s%d", carrier[:2], 100+d.seed.Intn(900)),
			Departure:  depTime,
			Arrival:    arrTime,
			DurationMin: durMin,
			Price:      price,
			Currency:   "EUR",
			BookingURL: fmt.Sprintf("https://www.google.com/travel/flights?q=flights+%s+to+%s", req.Origin, req.Destination),
			Cabin:      "ECONOMY",
		}

		trips = append(trips, models.Trip{
			ID:          fmt.Sprintf("demo-%d", i),
			Segments:    []models.Segment{seg},
			TotalPrice:  price,
			Currency:    "EUR",
			TotalDurMin: durMin,
			Source:      "demo",
			BookingURL:  seg.BookingURL,
			Transport:   []string{transportType},
		})
	}

	return trips, nil
}