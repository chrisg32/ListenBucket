package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/listenbucket/listenbucket/internal/database"
	"github.com/listenbucket/listenbucket/internal/podcast"
)

type CreateFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (s *Server) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := s.db.GetFeeds()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if feeds == nil {
		feeds = []database.Feed{}
	}

	json.NewEncoder(w).Encode(feeds)
}

func (s *Server) createFeed(w http.ResponseWriter, r *http.Request) {
	var req CreateFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	feed, err := s.db.CreateFeed(req.Title, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(feed)
}

func (s *Server) getFeed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	feed, err := s.db.GetFeed(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if feed == nil {
		http.Error(w, "feed not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(feed)
}

func (s *Server) updateFeed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.db.UpdateFeed(id, req.Title, req.Description); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteFeed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Check if it's the default feed
	feed, err := s.db.GetFeed(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if feed == nil {
		http.Error(w, "feed not found", http.StatusNotFound)
		return
	}

	if feed.IsDefault {
		http.Error(w, "cannot delete the default Listen Later feed", http.StatusForbidden)
		return
	}

	if err := s.db.DeleteFeed(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getFeedRSS(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	feed, err := s.db.GetFeed(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if feed == nil {
		http.Error(w, "feed not found", http.StatusNotFound)
		return
	}

	episodes, err := s.db.GetReadyEpisodesByFeed(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rss, err := podcast.GenerateFeed(feed, episodes, s.config.BaseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Write(rss)
}
