package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/listenbucket/listenbucket/internal/database"
)

// CreateSourceRequest represents the request body for adding a source
type CreateSourceRequest struct {
	URL string `json:"url"`
}

// listSources returns all sources for a feed
// GET /api/feeds/{feedId}/sources
func (s *Server) listSources(w http.ResponseWriter, r *http.Request) {
	feedID := chi.URLParam(r, "feedId")

	// Verify feed exists
	feed, err := s.db.GetFeed(feedID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}
	if feed == nil {
		writeError(w, http.StatusNotFound, "not_found", "feed not found")
		return
	}

	sources, err := s.db.GetSourcesByFeed(feedID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if sources == nil {
		sources = []database.Source{}
	}

	json.NewEncoder(w).Encode(sources)
}

// listAllSources returns all sources across all feeds
// GET /api/sources
func (s *Server) listAllSources(w http.ResponseWriter, r *http.Request) {
	sources, err := s.db.GetAllSources()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if sources == nil {
		sources = []database.Source{}
	}

	json.NewEncoder(w).Encode(sources)
}

// createSource adds a new source to a feed
// POST /api/feeds/{feedId}/sources
func (s *Server) createSource(w http.ResponseWriter, r *http.Request) {
	feedID := chi.URLParam(r, "feedId")

	// Verify feed exists
	feed, err := s.db.GetFeed(feedID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}
	if feed == nil {
		writeError(w, http.StatusNotFound, "not_found", "feed not found")
		return
	}

	var req CreateSourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	if req.URL == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "url is required")
		return
	}

	source, err := s.downloader.AddSource(feedID, req.URL)
	if err != nil {
		writeError(w, http.StatusBadRequest, "source_error", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, source)
}

// getSource returns a single source by ID
// GET /api/sources/{sourceId}
func (s *Server) getSource(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "sourceId")

	source, err := s.db.GetSource(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if source == nil {
		writeError(w, http.StatusNotFound, "not_found", "source not found")
		return
	}

	json.NewEncoder(w).Encode(source)
}

// deleteSource deletes a source
// DELETE /api/sources/{sourceId}
func (s *Server) deleteSource(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "sourceId")

	source, err := s.db.GetSource(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if source == nil {
		writeError(w, http.StatusNotFound, "not_found", "source not found")
		return
	}

	if err := s.db.DeleteSource(id); err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	writeSuccess(w, "source deleted successfully")
}

// refreshSource triggers a check for new content from a source
// POST /api/sources/{sourceId}/refresh
func (s *Server) refreshSource(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "sourceId")

	source, err := s.db.GetSource(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if source == nil {
		writeError(w, http.StatusNotFound, "not_found", "source not found")
		return
	}

	// Trigger refresh in background
	go s.downloader.RefreshSource(source)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "source refresh started",
		"source":  source,
	})
}
