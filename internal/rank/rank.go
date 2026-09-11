package rank

import (
	"math"
	"sort"
	"time"

	"github.com/dOuReallyDo/wanderer/internal/models"
)

// Engine calcola il Wander Score e ordina i trip.
type Engine struct {
	weights models.RankWeights
}

func NewEngine(w models.RankWeights) *Engine {
	return &Engine{weights: w}
}

// Rank applica il Wander Score a tutti i trip e li ordina.
func (e *Engine) Rank(trips []models.Trip, req models.SearchRequest) []models.Trip {
	if len(trips) == 0 {
		return trips
	}

	// Filtra per vincoli
	filtered := filter(trips, req)
	if len(filtered) == 0 {
		return filtered
	}

	// Normalizza metriche per confronto
	minPrice, maxPrice := boundsPrice(filtered)
	minDur, maxDur := boundsDur(filtered)

	for i := range filtered {
		t := &filtered[i]
		score := e.computeScore(t, req, minPrice, maxPrice, minDur, maxDur)
		t.Score = math.Round(score*100) / 100
	}

	// Ordina per score decrescente
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Score > filtered[j].Score
	})

	// Assegna rank
	for i := range filtered {
		filtered[i].Rank = i + 1
	}

	return filtered
}

func (e *Engine) computeScore(t *models.Trip, req models.SearchRequest,
	minP, maxP, minD, maxD float64) float64 {

	totalWeight := e.weights.Cost + e.weights.Duration + e.weights.Schedule + e.weights.Serendipity
	if totalWeight == 0 {
		totalWeight = 100
	}

	// Cost score: 100 se prezzo == min, 0 se == max
	costScore := normInverse(t.TotalPrice, minP, maxP)

	// Duration score: 100 se durata == min, 0 se == max
	durScore := normInverse(float64(t.TotalDurMin), minD, maxD)

	// Schedule score: preferisce partenze in orario comodo (8-20)
	scheduleScore := computeScheduleScore(t, req)

	// Serendipity: offerte con prezzo molto sotto la media
	serendipityScore := computeSerendipity(t, minP, maxP)

	raw := costScore*e.weights.Cost + durScore*e.weights.Duration +
		scheduleScore*e.weights.Schedule + serendipityScore*e.weights.Serendipity

	return (raw / totalWeight)
}

// normInverse: 100 se val==min, 0 se val==max (più basso = meglio)
func normInverse(val, minV, maxV float64) float64 {
	if maxV == minV {
		return 100
	}
	return 100 * (maxV - val) / (maxV - minV)
}

func computeScheduleScore(t *models.Trip, req models.SearchRequest) float64 {
	if len(t.Segments) == 0 {
		return 50
	}
	dep := t.Segments[0].Departure
	hour := dep.Hour()

	// Vincoli espliciti
	if req.EarliestDepart != "" && req.LatestDepart != "" {
		eh, _ := time.Parse("15:04", req.EarliestDepart)
		lh, _ := time.Parse("15:04", req.LatestDepart)
		depHM, _ := time.Parse("15:04", dep.Format("15:04"))
		if depHM.Before(eh) || depHM.After(lh) {
			return 0 // fuori vincolo
		}
		return 100
	}

	// Score morbido: 8-20 = ottimo, fuori penalizzato
	if hour >= 8 && hour <= 20 {
		return 100
	}
	if hour >= 6 && hour <= 22 {
		return 60
	}
	return 25
}

func computeSerendipity(t *models.Trip, minP, maxP float64) float64 {
	if maxP == minP {
		return 50
	}
	// Se prezzo nel bottom 25% della forchetta → high serendipity
	threshold := minP + (maxP-minP)*0.25
	if t.TotalPrice <= threshold {
		return 100
	}
	if t.TotalPrice <= minP+(maxP-minP)*0.5 {
		return 60
	}
	return 20
}

func filter(trips []models.Trip, req models.SearchRequest) []models.Trip {
	var out []models.Trip
	for _, t := range trips {
		// Filtro prezzo
		if req.MaxPrice > 0 && t.TotalPrice > req.MaxPrice {
			continue
		}
		// Filtro durata
		if req.MaxDurationH > 0 && t.TotalDurMin > req.MaxDurationH*60 {
			continue
		}
		// Filtro tipo transporte
		if req.PreferredTransport != "" && req.PreferredTransport != "any" {
			match := false
			for _, tr := range t.Transport {
				if tr == req.PreferredTransport {
					match = true
				}
			}
			if !match && req.PreferredTransport == "mixed" {
				match = len(t.Transport) > 1
			}
			if !match {
				continue
			}
		}
		out = append(out, t)
	}
	return out
}

func boundsPrice(trips []models.Trip) (min, max float64) {
	min, max = math.MaxFloat64, 0
	for _, t := range trips {
		if t.TotalPrice < min {
			min = t.TotalPrice
		}
		if t.TotalPrice > max {
			max = t.TotalPrice
		}
	}
	return
}

func boundsDur(trips []models.Trip) (min, max float64) {
	min, max = math.MaxFloat64, 0
	for _, t := range trips {
		d := float64(t.TotalDurMin)
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}
	return
}