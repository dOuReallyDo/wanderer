package search

import (
	"context"
	"time"

	"github.com/dOuReallyDo/wanderer/internal/models"
	"github.com/dOuReallyDo/wanderer/internal/providers"
	"github.com/dOuReallyDo/wanderer/internal/rank"
)

// Engine orchestra la ricerca cross-provider + ranking.
type Engine struct {
	registry *providers.Registry
	ranker   *rank.Engine
	timeout  time.Duration
}

func NewEngine(reg *providers.Registry, r *rank.Engine) *Engine {
	return &Engine{
		registry: reg,
		ranker:   r,
		timeout:  20 * time.Second,
	}
}

func (e *Engine) Search(req models.SearchRequest) models.SearchResponse {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	trips, queried := e.registry.SearchAll(ctx, req)
	ranked := e.ranker.Rank(trips, req)

	resp := models.SearchResponse{
		Request:    req,
		Trips:      ranked,
		Count:      len(ranked),
		SearchedAt: start,
		Providers:  queried,
	}
	if len(ranked) > 0 {
		resp.Best = &ranked[0]
	}
	resp.TookMs = time.Since(start).Milliseconds()
	return resp
}