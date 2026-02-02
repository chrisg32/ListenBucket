package api

import (
	"embed"
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

		// Feeds
		r.Get("/feeds", s.getFeeds)
		r.Post("/feeds", s.createFeed)
		r.Get("/feeds/{id}", s.getFeed)
		r.Put("/feeds/{id}", s.updateFeed)
		r.Delete("/feeds/{id}", s.deleteFeed)
		r.Get("/feeds/{id}/rss", s.getFeedRSS)

		// Episodes
		r.Get("/feeds/{feedId}/episodes", s.getEpisodes)
		r.Delete("/episodes/{id}", s.deleteEpisode)

		// Sources
		r.Get("/feeds/{feedId}/sources", s.getSources)
		r.Post("/feeds/{feedId}/sources", s.addSource)
		r.Delete("/sources/{id}", s.deleteSource)

		// Media (no JSON content type)
		r.Group(func(r chi.Router) {
			r.Get("/media/{filename}", s.serveMedia)
		})
	})

	// Serve static files from embedded filesystem
	distFS, err := fs.Sub(webFS, "web/dist")
	if err == nil {
		fileServer := http.FileServer(http.FS(distFS))
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			// Try to serve the file
			path := r.URL.Path
			if path == "/" {
				path = "/index.html"
			}

			// Check if file exists
			f, err := distFS.Open(filepath.Clean(path[1:]))
			if err != nil {
				// Serve index.html for SPA routing
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
