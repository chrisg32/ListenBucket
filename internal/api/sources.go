package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/listenbucket/listenbucket/internal/database"
)

type AddSourceRequest struct {
	URL string `json:"url"`
}

func (s *Server) getSources(w http.ResponseWriter, r *http.Request) {
	feedID := chi.URLParam(r, "feedId")

	sources, err := s.db.GetSourcesByFeed(feedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if sources == nil {
		sources = []database.Source{}
	}

	json.NewEncoder(w).Encode(sources)
}

func (s *Server) addSource(w http.ResponseWriter, r *http.Request) {
	feedID := chi.URLParam(r, "feedId")

	var req AddSourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	// Verify feed exists
	feed, err := s.db.GetFeed(feedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if feed == nil {
		http.Error(w, "feed not found", http.StatusNotFound)
		return
	}

	source, err := s.downloader.AddSource(feedID, req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(source)
}

func (s *Server) deleteSource(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := s.db.DeleteSource(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
