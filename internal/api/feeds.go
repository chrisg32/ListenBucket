package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/chrisg32/ListenBucket/internal/database"
	"github.com/chrisg32/ListenBucket/internal/podcast"
)

// CreateFeedRequest represents the request body for creating a feed
type CreateFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateFeedRequest represents the request body for updating a feed
type UpdateFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// listFeeds returns all feeds
// GET /api/feeds
func (s *Server) listFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := s.db.GetFeeds()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if feeds == nil {
		feeds = []database.Feed{}
	}

	json.NewEncoder(w).Encode(feeds)
}

// createFeed creates a new feed
// POST /api/feeds
func (s *Server) createFeed(w http.ResponseWriter, r *http.Request) {
	var req CreateFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "title is required")
		return
	}

	feed, err := s.db.CreateFeed(req.Title, req.Description)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, feed)
}

// getFeed returns a single feed by ID
// GET /api/feeds/{feedId}
func (s *Server) getFeed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "feedId")

	feed, err := s.db.GetFeed(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if feed == nil {
		writeError(w, http.StatusNotFound, "not_found", "feed not found")
		return
	}

	json.NewEncoder(w).Encode(feed)
}

// updateFeed updates an existing feed
// PUT /api/feeds/{feedId}
func (s *Server) updateFeed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "feedId")

	// Check if feed exists
	feed, err := s.db.GetFeed(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}
	if feed == nil {
		writeError(w, http.StatusNotFound, "not_found", "feed not found")
		return
	}

	var req UpdateFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "title is required")
		return
	}

	if err := s.db.UpdateFeed(id, req.Title, req.Description); err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	// Return updated feed
	updatedFeed, _ := s.db.GetFeed(id)
	json.NewEncoder(w).Encode(updatedFeed)
}

// deleteFeed deletes a feed
// DELETE /api/feeds/{feedId}
func (s *Server) deleteFeed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "feedId")

	feed, err := s.db.GetFeed(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if feed == nil {
		writeError(w, http.StatusNotFound, "not_found", "feed not found")
		return
	}

	if feed.IsDefault {
		writeError(w, http.StatusForbidden, "forbidden", "cannot delete the default Listen Later feed")
		return
	}

	if err := s.db.DeleteFeed(id); err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	writeSuccess(w, "feed deleted successfully")
}

// getFeedRSS returns the RSS feed XML
// GET /api/feeds/{feedId}/rss
func (s *Server) getFeedRSS(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "feedId")

	feed, err := s.db.GetFeed(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if feed == nil {
		writeError(w, http.StatusNotFound, "not_found", "feed not found")
		return
	}

	episodes, err := s.db.GetReadyEpisodesByFeed(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	// Determine base URL from request to ensure URLs work from any device
	baseURL := s.config.BaseURL
	if r.Host != "" {
		scheme := "http"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}
		baseURL = scheme + "://" + r.Host
	}

	rss, err := podcast.GenerateFeed(feed, episodes, baseURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "rss_generation_error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Write(rss)
}
