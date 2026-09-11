package providers

import (
	"context"
	"sync"

	"github.com/dOuReallyDo/wanderer/internal/models"
)

// Provider è l'interfaccia che ogni sorgente di offerte implementa.
type Provider interface {
	Name() string
	Search(ctx context.Context, req models.SearchRequest) ([]models.Trip, error)
	Available() bool
}

// Registry gestisce i provider registrati.
type Registry struct {
	mu        sync.RWMutex
	providers []Provider
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) Register(p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = append(r.providers, p)
}

func (r *Registry) All() []Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cp := make([]Provider, len(r.providers))
	copy(cp, r.providers)
	return cp
}

// SearchAll interroga tutti i provider attivi in parallelo.
func (r *Registry) SearchAll(ctx context.Context, req models.SearchRequest) ([]models.Trip, []string) {
	providers := r.All()
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allTrips []models.Trip
	var queried []string

	for _, p := range providers {
		if !p.Available() {
			continue
		}
		wg.Add(1)
		queried = append(queried, p.Name())
		go func(p Provider) {
			defer wg.Done()
			trips, err := p.Search(ctx, req)
			if err != nil {
				return // provider fallito: saltiamo, altri continuano
			}
			mu.Lock()
			allTrips = append(allTrips, trips...)
			mu.Unlock()
		}(p)
	}
	wg.Wait()
	return allTrips, queried
}