package http

import (
	"github.com/Quickdesh/SyncMiru/internal/database"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type healthHandler struct {
	encoder encoder
	db      *database.DB
}

func newHealthHandler(encoder encoder, db *database.DB) *healthHandler {
	return &healthHandler{
		encoder: encoder,
		db:      db,
	}
}

func (h healthHandler) Routes(r chi.Router) {
	r.Get("/liveness", h.handleLiveness)
	r.Get("/readiness", h.handleReadiness)
}

func (h healthHandler) handleLiveness(w http.ResponseWriter, _ *http.Request) {
	writeHealthy(w)
}

func (h healthHandler) handleReadiness(w http.ResponseWriter, _ *http.Request) {
	if err := h.db.Ping(); err != nil {
		writeUnhealthy(w)
		return
	}

	writeHealthy(w)
}

func writeHealthy(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func writeUnhealthy(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Unhealthy. Database unreachable"))
}
