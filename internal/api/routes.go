package api

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/chrisg32/ListenBucket/internal/assets"
	"github.com/chrisg32/ListenBucket/internal/config"
	"github.com/chrisg32/ListenBucket/internal/database"
	"github.com/chrisg32/ListenBucket/internal/downloader"
)

type Server struct {
	db         *database.DB
	downloader *downloader.Downloader
	config     *config.Config
}

// APIError represents a JSON error response
type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// APISuccess represents a JSON success response for operations without data
type APISuccess struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func NewServer(db *database.DB, dl *downloader.Downloader, cfg *config.Config) *Server {
	return &Server{
		db:         db,
		downloader: dl,
		config:     cfg,
	}
}

func (s *Server) SetupRoutes(webFS embed.FS) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Use(jsonContentType)

		// Health check (public)
		r.Get("/health", s.healthCheck)

		// Auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Get("/setup-status", s.getSetupStatus)
			r.Post("/setup", s.setup)
			r.Post("/login", s.login)
			r.Post("/logout", s.logout)
			r.With(s.authMiddleware).Get("/me", s.me)
		})

		// RSS feed (public - for podcast apps)
		r.Get("/feeds/{feedId}/rss", s.getFeedRSS)

		// Media files (public - for podcast apps)
		r.Get("/media/{filename}", s.serveMedia)

		// Logo (public)
		r.Get("/logo.png", s.serveLogo)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(s.authMiddleware)

			// Feeds - RESTful CRUD
			r.Get("/feeds", s.listFeeds)
			r.Post("/feeds", s.createFeed)
			r.Get("/feeds/{feedId}", s.getFeed)
			r.Put("/feeds/{feedId}", s.updateFeed)
			r.Delete("/feeds/{feedId}", s.deleteFeed)

			// Episodes nested under feeds
			r.Get("/feeds/{feedId}/episodes", s.listEpisodes)
			r.Post("/feeds/{feedId}/episodes", s.createEpisode)
			r.Post("/feeds/{feedId}/episodes/upload", s.uploadEpisode)

			// Sources nested under feeds
			r.Get("/feeds/{feedId}/sources", s.listSources)
			r.Post("/feeds/{feedId}/sources", s.createSource)

			// Episodes - direct access by ID
			r.Get("/episodes", s.listAllEpisodes)
			r.Get("/episodes/{episodeId}", s.getEpisode)
			r.Put("/episodes/{episodeId}", s.updateEpisode)
			r.Delete("/episodes/{episodeId}", s.deleteEpisode)
			r.Post("/episodes/{episodeId}/retry", s.retryEpisode)

			// Sources - direct access by ID
			r.Get("/sources", s.listAllSources)
			r.Get("/sources/{sourceId}", s.getSource)
			r.Delete("/sources/{sourceId}", s.deleteSource)
			r.Post("/sources/{sourceId}/refresh", s.refreshSource)
		})
	})

	// Serve static files from embedded filesystem
	distFS, err := fs.Sub(webFS, "web/dist")
	if err == nil {
		fileServer := http.FileServer(http.FS(distFS))
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if path == "/" {
				path = "/index.html"
			}

			f, err := distFS.Open(filepath.Clean(path[1:]))
			if err != nil {
				r.URL.Path = "/"
				fileServer.ServeHTTP(w, r)
				return
			}
			f.Close()

			fileServer.ServeHTTP(w, r)
		})
	}

	return r
}

func jsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// writeError writes a JSON error response
func writeError(w http.ResponseWriter, code int, err string, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(APIError{
		Error:   err,
		Message: message,
		Code:    code,
	})
}

// writeSuccess writes a JSON success response
func writeSuccess(w http.ResponseWriter, message string) {
	json.NewEncoder(w).Encode(APISuccess{
		Success: true,
		Message: message,
	})
}

// writeJSON writes a JSON response with the given status code
func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

// healthCheck returns the health status of the API
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"version": "1.1.0",
	})
}

// serveLogo serves the ListenBucket logo
func (s *Server) serveLogo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Write(assets.Logo)
}
