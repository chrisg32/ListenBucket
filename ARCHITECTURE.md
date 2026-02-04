# ListenBucket Architecture

This document describes the architecture and design decisions for ListenBucket. It is intended as a reference for AI agents and developers working on the codebase.

## Overview

ListenBucket is a web application that converts YouTube videos, playlists, and channels into podcast feeds. Users can subscribe to these feeds in any podcast app and listen to YouTube content as audio.

## Design Philosophy

### Ease of Deployment
- **Single binary architecture** - The entire application (backend + frontend) compiles to a single executable
- **Minimal dependencies** - Only requires yt-dlp and ffmpeg at runtime
- **Docker-first** - Optimized for containerized deployment
- **CasaOS compatible** - Easy installation on home servers

### Performance
- **Low memory footprint** - Suitable for low-power devices (Raspberry Pi, NAS)
- **Efficient audio extraction** - Uses yt-dlp for downloading and ffmpeg for conversion
- **SQLite database** - Simple, file-based storage with no external database server

## Technology Stack

### Backend: Go (Golang)

**Why Go?**
- Compiles to a single static binary
- Excellent performance with low resource usage
- Strong standard library for HTTP servers
- `go:embed` allows embedding the frontend in the binary
- Patterns are standardized, making it easy for LLMs to write

**Key Libraries:**
- `github.com/go-chi/chi/v5` - HTTP router
- `modernc.org/sqlite` - Pure Go SQLite (no CGO required)
- `github.com/google/uuid` - UUID generation

### Frontend: Vue.js 3 + Tailwind CSS

**Why Vue.js?**
- Lightweight and fast
- Component-based architecture
- Excellent tooling with Vite
- Similar to Mailpit's architecture

**Why Tailwind CSS?**
- Utility-first approach works well with AI code generation
- Consistent styling without custom CSS files
- Dark mode support built-in

**Key Libraries:**
- `vue-router` - Client-side routing
- `pinia` - State management

### Data Storage: SQLite

**Why SQLite?**
- No external database server required
- File-based, easy to backup
- `modernc.org/sqlite` provides pure Go implementation (no CGO)
- WAL mode for better concurrency

**Schema:**
- `feeds` - Podcast feeds (default "Listen Later" + user-created)
- `episodes` - Individual podcast episodes
- `sources` - YouTube videos/playlists/channels linked to feeds

**Migrations:**
- Uses `golang-migrate/migrate` for versioned schema migrations
- Migration files in `internal/database/migrations/` (embedded via `go:embed`)
- Migrations run automatically on startup
- Version tracking in `schema_migrations` table
- Supports both up and down migrations for rollbacks

## Project Structure

```
.
├── cmd/listenbucket/          # Application entry point
│   └── main.go
├── internal/
│   ├── api/                   # HTTP handlers and routing
│   │   ├── routes.go          # Route definitions
│   │   ├── feeds.go           # Feed CRUD handlers
│   │   ├── episodes.go        # Episode handlers
│   │   └── sources.go         # Source handlers
│   ├── assets/                # Embedded assets (logo via go:embed)
│   ├── config/                # Configuration loading
│   ├── database/              # SQLite database layer
│   │   ├── db.go              # Database operations
│   │   ├── models.go          # Data models
│   │   ├── migrations.go      # Embedded migrations
│   │   └── migrations/        # SQL migration files
│   ├── downloader/            # yt-dlp integration
│   │   └── ytdlp.go           # Download and extraction logic
│   └── podcast/               # RSS feed generation
│       └── rss.go             # iTunes-compatible RSS
├── web/                       # Vue.js frontend
│   ├── src/
│   │   ├── components/        # Reusable Vue components
│   │   ├── views/             # Page components
│   │   ├── stores/            # Pinia state stores
│   │   └── router/            # Vue Router config
│   ├── index.html
│   └── package.json
├── Dockerfile                 # Multi-stage Docker build
├── docker-compose.yml
├── docker-compose.casaos.yml  # CasaOS-specific config
├── Makefile
└── .github/workflows/         # CI/CD automation
    ├── ci.yml                 # Test, lint, and build
    └── release.yml            # Docker image publishing
```

## Key Workflows

### Adding a Source

1. User submits YouTube URL via frontend
2. User optionally enables "Include back catalog" for playlists/channels
3. Backend detects source type (video/playlist/channel)
4. yt-dlp fetches metadata (title, thumbnail, video list)
5. Source record created in database
6. For videos: Episode created immediately
7. For playlists/channels:
   - If `include_back_catalog` is true: All existing videos added as episodes
   - If false: Only new videos added going forward

### Episode Download

1. Background worker polls for pending episodes every 30 seconds
2. For each pending episode:
   - Status updated to "downloading"
   - yt-dlp downloads and extracts audio
   - ffmpeg converts to MP3
   - Duration calculated with ffprobe
   - Status updated to "ready" with audio URL

### RSS Feed Generation

1. Request to `/api/feeds/{id}/rss`
2. Fetch feed metadata and ready episodes
3. Generate iTunes-compatible RSS XML
4. Audio URLs constructed using request host (for remote access)

### Source Checking

1. Background worker runs every hour
2. Fetches playlists/channels that need updating
3. Compares current videos with existing episodes
4. Creates episodes for new videos

## API Design

RESTful API with JSON responses:
- `GET/POST /api/feeds` - List/create feeds
- `GET/PUT/DELETE /api/feeds/{id}` - Feed CRUD
- `GET /api/feeds/{id}/rss` - RSS feed
- `GET/POST /api/feeds/{id}/episodes` - Feed episodes
- `POST /api/feeds/{id}/episodes/upload` - Upload audio/video file (multipart form)
- `GET/POST /api/feeds/{id}/sources` - Feed sources
- `GET/PUT/DELETE /api/episodes/{id}` - Episode CRUD
- `GET/DELETE /api/sources/{id}` - Source CRUD

## Frontend Architecture

### State Management (Pinia)
- `feeds` store manages all feed, episode, and source data
- `auth` store manages authentication state and user info
- API calls centralized in store actions
- Reactive updates when data changes

### Routing (Vue Router)
- `/` - Home (feed list) - requires auth
- `/feed/:id` - Feed details with tabs for episodes/sources - requires auth
- `/login` - Login page
- `/setup` - Initial setup/onboarding page

### Components
- `FeedCard` - Feed display on home page
- `EpisodeCard` - Episode with embedded audio player
- `SourceCard` - Source with type badge (video/playlist/channel)
- `AudioPlayer` - Custom themed audio player with seek, progress bar, and time display
- `AddSourceModal` - Modal for adding YouTube URLs with "include back catalog" option
- `UploadFileModal` - Modal for uploading audio/video files with drag-and-drop
- `CreateFeedModal` - Modal for creating new feeds
- `FeedLinks` - RSS/Podcast app subscription links
- `Header` - Navigation bar with theme toggle

## Docker Build

Multi-stage Dockerfile:
1. **frontend-builder** - Node.js builds Vue.js to `dist/`
2. **backend-builder** - Go compiles with embedded frontend
3. **runtime** - Alpine with yt-dlp and ffmpeg

## CI/CD

GitHub Actions workflows in `.github/workflows/`:

### CI Pipeline (`ci.yml`)
Runs on push/PR to `main`:
1. **Test** - Go tests with race detection and coverage (uploaded to Codecov)
2. **Lint** - golangci-lint for code quality
3. **Build Docker** - Validates Docker image builds (depends on test/lint)

### Release Pipeline (`release.yml`)
Publishes Docker images to Docker Hub on releases.

## Configuration

Environment variables:
- `PORT` - HTTP server port (default: 8080)
- `DATA_DIR` - Database and media storage (default: ./data)
- `BASE_URL` - Public URL for RSS feeds (default: http://localhost:8080)

## Authentication

Session-based authentication with bcrypt password hashing:
- Initial setup creates the first user account
- Cookie-based sessions (30-day expiration)
- Protected routes require authentication
- RSS feeds and media files remain public for podcast apps

**Auth Endpoints:**
- `GET /api/auth/setup-status` - Check if initial setup is required
- `POST /api/auth/setup` - Create initial user account
- `POST /api/auth/login` - Authenticate user
- `POST /api/auth/logout` - Invalidate session
- `GET /api/auth/me` - Get current user

**Database Tables:**
- `users` - User accounts with bcrypt-hashed passwords
- `sessions` - Active sessions with expiration tracking

## Future Considerations

Features planned but not yet implemented:
- S3/cloud storage for media files

## Development Guidelines

1. **Keep it simple** - Avoid over-engineering
2. **Single responsibility** - Each package has a clear purpose
3. **Error handling** - Always return meaningful errors
4. **Testing** - Write tests for database and API layers
5. **Documentation** - Update this file when architecture changes
