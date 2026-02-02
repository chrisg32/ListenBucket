package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	"github.com/listenbucket/listenbucket/internal/database"
)

func (s *Server) getEpisodes(w http.ResponseWriter, r *http.Request) {
	feedID := chi.URLParam(r, "feedId")

	episodes, err := s.db.GetEpisodesByFeed(feedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if episodes == nil {
		episodes = []database.Episode{}
	}

	json.NewEncoder(w).Encode(episodes)
}

func (s *Server) deleteEpisode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Get episode to find the audio file
	ep, err := s.db.GetEpisode(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ep != nil {
		// Try to delete the audio file
		audioPath := filepath.Join(s.config.MediaDir, id+".mp3")
		os.Remove(audioPath) // Ignore error if file doesn't exist
	}

	if err := s.db.DeleteEpisode(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) serveMedia(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")

	// Security: ensure filename doesn't contain path traversal
	if filepath.Base(filename) != filename {
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(s.config.MediaDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	http.ServeFile(w, r, filePath)
}
