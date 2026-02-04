package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	"github.com/chrisg32/ListenBucket/internal/database"
)

// CreateEpisodeRequest represents the request body for creating an episode
type CreateEpisodeRequest struct {
	URL         string `json:"url"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateEpisodeRequest represents the request body for updating an episode
type UpdateEpisodeRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// listEpisodes returns all episodes for a feed
// GET /api/feeds/{feedId}/episodes
func (s *Server) listEpisodes(w http.ResponseWriter, r *http.Request) {
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

	episodes, err := s.db.GetEpisodesByFeed(feedID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if episodes == nil {
		episodes = []database.Episode{}
	}

	json.NewEncoder(w).Encode(episodes)
}

// listAllEpisodes returns all episodes across all feeds
// GET /api/episodes
func (s *Server) listAllEpisodes(w http.ResponseWriter, r *http.Request) {
	episodes, err := s.db.GetAllEpisodes()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if episodes == nil {
		episodes = []database.Episode{}
	}

	json.NewEncoder(w).Encode(episodes)
}

// createEpisode adds a new episode to a feed from a URL
// POST /api/feeds/{feedId}/episodes
func (s *Server) createEpisode(w http.ResponseWriter, r *http.Request) {
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

	var req CreateEpisodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	if req.URL == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "url is required")
		return
	}

	// Get video info from URL
	info, err := s.downloader.GetVideoInfo(req.URL)
	if err != nil {
		writeError(w, http.StatusBadRequest, "url_error", "failed to get video info: "+err.Error())
		return
	}

	title := req.Title
	if title == "" {
		title = info.Title
	}

	description := req.Description
	if description == "" {
		description = info.Description
	}

	episode, err := s.db.CreateEpisode(feedID, "", title, description, info.Thumbnail, req.URL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, episode)
}

// getEpisode returns a single episode by ID
// GET /api/episodes/{episodeId}
func (s *Server) getEpisode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "episodeId")

	episode, err := s.db.GetEpisode(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if episode == nil {
		writeError(w, http.StatusNotFound, "not_found", "episode not found")
		return
	}

	json.NewEncoder(w).Encode(episode)
}

// updateEpisode updates an episode's metadata
// PUT /api/episodes/{episodeId}
func (s *Server) updateEpisode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "episodeId")

	episode, err := s.db.GetEpisode(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if episode == nil {
		writeError(w, http.StatusNotFound, "not_found", "episode not found")
		return
	}

	var req UpdateEpisodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	// Update fields if provided
	title := episode.Title
	description := episode.Description
	if req.Title != "" {
		title = req.Title
	}
	if req.Description != "" {
		description = req.Description
	}

	if err := s.db.UpdateEpisodeMetadata(id, title, description); err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	// Return updated episode
	updated, _ := s.db.GetEpisode(id)
	json.NewEncoder(w).Encode(updated)
}

// deleteEpisode deletes an episode
// DELETE /api/episodes/{episodeId}
func (s *Server) deleteEpisode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "episodeId")

	episode, err := s.db.GetEpisode(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if episode == nil {
		writeError(w, http.StatusNotFound, "not_found", "episode not found")
		return
	}

	// Delete the audio file if it exists
	audioPath := filepath.Join(s.config.MediaDir, id+".mp3")
	os.Remove(audioPath)

	if err := s.db.DeleteEpisode(id); err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	writeSuccess(w, "episode deleted successfully")
}

// retryEpisode resets an episode's status to retry download
// POST /api/episodes/{episodeId}/retry
func (s *Server) retryEpisode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "episodeId")

	episode, err := s.db.GetEpisode(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	if episode == nil {
		writeError(w, http.StatusNotFound, "not_found", "episode not found")
		return
	}

	// Only allow retry for error or pending status
	if episode.Status != database.StatusError && episode.Status != database.StatusPending {
		writeError(w, http.StatusBadRequest, "invalid_status", "can only retry episodes with error or pending status")
		return
	}

	// Reset status to pending
	if err := s.db.UpdateEpisodeStatus(id, database.StatusPending, ""); err != nil {
		writeError(w, http.StatusInternalServerError, "database_error", err.Error())
		return
	}

	// Return updated episode
	updated, _ := s.db.GetEpisode(id)
	json.NewEncoder(w).Encode(updated)
}

// serveMedia serves media files
// GET /api/media/{filename}
func (s *Server) serveMedia(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")

	// Security: ensure filename doesn't contain path traversal
	if filepath.Base(filename) != filename {
		writeError(w, http.StatusBadRequest, "invalid_filename", "invalid filename")
		return
	}

	filePath := filepath.Join(s.config.MediaDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		writeError(w, http.StatusNotFound, "not_found", "file not found")
		return
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	http.ServeFile(w, r, filePath)
}
