package api

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/listenbucket/listenbucket/internal/config"
	"github.com/listenbucket/listenbucket/internal/database"
	"github.com/listenbucket/listenbucket/internal/downloader"
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

		// Health check
		r.Get("/health", s.healthCheck)

		// Feeds - RESTful CRUD
		r.Get("/feeds", s.listFeeds)             // GET /api/feeds - List all feeds
		r.Post("/feeds", s.createFeed)           // POST /api/feeds - Create new feed
		r.Get("/feeds/{feedId}", s.getFeed)      // GET /api/feeds/{feedId} - Get feed
		r.Put("/feeds/{feedId}", s.updateFeed)   // PUT /api/feeds/{feedId} - Update feed
		r.Delete("/feeds/{feedId}", s.deleteFeed) // DELETE /api/feeds/{feedId} - Delete feed
		r.Get("/feeds/{feedId}/rss", s.getFeedRSS) // GET /api/feeds/{feedId}/rss - Get RSS XML

		// Episodes nested under feeds
		r.Get("/feeds/{feedId}/episodes", s.listEpisodes)   // GET /api/feeds/{feedId}/episodes
		r.Post("/feeds/{feedId}/episodes", s.createEpisode) // POST /api/feeds/{feedId}/episodes

		// Sources nested under feeds
		r.Get("/feeds/{feedId}/sources", s.listSources)   // GET /api/feeds/{feedId}/sources
		r.Post("/feeds/{feedId}/sources", s.createSource) // POST /api/feeds/{feedId}/sources

		// Episodes - direct access by ID
		r.Get("/episodes", s.listAllEpisodes)              // GET /api/episodes - List all episodes
		r.Get("/episodes/{episodeId}", s.getEpisode)       // GET /api/episodes/{episodeId}
		r.Put("/episodes/{episodeId}", s.updateEpisode)    // PUT /api/episodes/{episodeId}
		r.Delete("/episodes/{episodeId}", s.deleteEpisode) // DELETE /api/episodes/{episodeId}
		r.Post("/episodes/{episodeId}/retry", s.retryEpisode) // POST /api/episodes/{episodeId}/retry

		// Sources - direct access by ID
		r.Get("/sources", s.listAllSources)                // GET /api/sources - List all sources
		r.Get("/sources/{sourceId}", s.getSource)          // GET /api/sources/{sourceId}
		r.Delete("/sources/{sourceId}", s.deleteSource)    // DELETE /api/sources/{sourceId}
		r.Post("/sources/{sourceId}/refresh", s.refreshSource) // POST /api/sources/{sourceId}/refresh

		// Media files (different content type)
		r.Get("/media/{filename}", s.serveMedia)
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
		"version": "1.0.0",
	})
}
