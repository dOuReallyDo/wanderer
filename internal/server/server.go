package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dOuReallyDo/wanderer/internal/models"
	"github.com/dOuReallyDo/wanderer/internal/search"
)

type Server struct {
	engine *search.Engine
	mux    *http.ServeMux
	port   string
}

func New(port string, engine *search.Engine) *Server {
	s := &Server{
		engine: engine,
		mux:    http.NewServeMux(),
		port:   port,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	// API
	s.mux.HandleFunc("/api/search", s.handleSearch)
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/weights", s.handleWeights)
	s.mux.HandleFunc("/api/cities", s.handleCities)

	// Static UI (embedded)
	s.mux.Handle("/", http.FileServerFS(staticFS))
}

func (s *Server) Listen() error {
	addr := ":" + s.port
	fmt.Fprintf(os.Stderr, "🌐 Wanderer su http://localhost:%s\n", s.port)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(204)
		return
	}

	var req models.SearchRequest
	if r.Method == "GET" {
		req = models.SearchRequest{
			Origin:       r.URL.Query().Get("origin"),
			Destination:  r.URL.Query().Get("destination"),
			Date:         r.URL.Query().Get("date"),
			ReturnDate:   r.URL.Query().Get("return_date"),
			Adults:       parseIntDefault(r.URL.Query().Get("adults"), 1),
			MaxPrice:     parseFloat(r.URL.Query().Get("max_price")),
			Currency:     strOr(r.URL.Query().Get("currency"), "EUR"),
			MaxDurationH: parseIntDefault(r.URL.Query().Get("max_duration_h"), 0),
			EarliestDepart: r.URL.Query().Get("earliest_depart"),
			LatestDepart:   r.URL.Query().Get("latest_depart"),
			PreferredTransport: r.URL.Query().Get("preferred_transport"),
		}
	} else if r.Method == "POST" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, 400, "invalid JSON: "+err.Error())
			return
		}
	} else {
		writeError(w, 405, "method not allowed")
		return
	}

	if req.Origin == "" || req.Destination == "" || req.Date == "" {
		writeError(w, 400, "origin, destination, date sono obbligatori")
		return
	}

	resp := s.engine.Search(req)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (s *Server) handleWeights(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "GET" {
		json.NewEncoder(w).Encode(models.DefaultWeights())
		return
	}

	var w2 models.RankWeights
	if err := json.NewDecoder(r.Body).Decode(&w2); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	json.NewEncoder(w).Encode(w2)
}

func (s *Server) handleCities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	q := r.URL.Query().Get("q")
	cities := models.SearchCities(q)
	json.NewEncoder(w).Encode(map[string]any{
		"query":  q,
		"count":  len(cities),
		"cities": cities,
	})
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	var v int
	fmt.Sscanf(s, "%d", &v)
	if v == 0 {
		return def
	}
	return v
}

func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	var v float64
	fmt.Sscanf(s, "%f", &v)
	return v
}

func strOr(s, def string) string {
	if s == "" {
		return def
	}
	return s
}