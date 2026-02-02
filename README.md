# ListenBucket
 
### Project Description:

ListenBucket is a web application that allows users to insert a link to a video or audio file. It will the add that file as an episode to a podcast feed that the application hosts capturing any meta-data from the source and encoding/extracting the audio from the video file.

There is one feed, called Listen Later, that every user has and cannot be deleted.

Besides listen later, users can create additional podcast feeds. Users can delete feeds that they create.

This is essentially what the ListenBox app (discontinued) did before YouTube shut it down.

Pages:
- Home page/list of feeds
	- Lists feeds (Image, Title, Description, last updated)
	- Can add/delete feeds (except ListenLater)
	- Can open feeds in podcast app
- Feed details
	- Primary Tab
		- Lists add the podcast episodes (Image, Title, Description, Data Added)
		- Ability to play episode audio inline
		- If no source only show the sources tab
	- Sources Tab
		- Ability to add sources
		- Ability to remove/delete sources
### Features:
- Free and open source
- Responsive UI, works on mobile and desktop
- Easy, docker install for self hosting
- Easy install into CasaOS
- YouTube video support
- YouTube playlist support - adds all videos from playlist, periodically checks for updates
- YouTube channel support - adds all videos from channel, periodically checks for updates
- From the feeds and feed details there are links to open the feed: direct rss link, podcasts app (ios), overcast app (ios)
- Episode image in feed is the thumbnail image for youtube (youtube hosted, just use url)
- Feed/podcast images come from the image of the source(s) (videos, playlists, or channels)
- Light and Dark UI Themes

### Future features (Do NOT implement yet):
- Directly upload a file or a video
- Multiple users
- Authentication/Authorization
- CANCELLED - Cloud mirroring - will upload feed contents (xml + mp3s) to a cloud provider and use the external cloud provider link for opening in podcast apps, easy way to get around firewalls
	- S3 bucket mirror

### Implementation guidelines:

Ease of deployment is critical. It should be able to be hosted anywhere with minimal dependencies.

It should be fast. It should be efficient.

It should utilize ffmpeg and/or yt-dlp.

All code should be thoroughly tested.

This stack mirrors the architecture of [Mailpit](https://mailpit.axllent.org/) (which uses Go + Vue.js embedded in a single binary). It offers the easiest deployment for a side project.

- **Front End:** **Vue.js + Tailwind CSS**
    
    - Mailpit uses Vue.js. Tailwind is easier for AI to write/refactor than custom CSS.
        
- **Back End:** **Go (Golang) + `lrstanley/go-ytdlp`**
    
    - **HTTP:** Go Standard Library `net/http` or `Chi` (keep it simple).
        
    - **yt-dlp:** `github.com/lrstanley/go-ytdlp` (Active, fluent API, supports context).
        
- **Reasons for Choice:**
    
    - **Deployment:** Go compiles to a **single static binary** that contains your entire frontend (via `//go:embed`). You copy _one file_ to your server/CasaOS, and it runs. No `node_modules`, no runtime dependencies.
        
    - **Vibe Coding:** Go is verbose but explicit. LLMs (Gemini/ChatGPT) are excellent at writing robust Go code because the patterns are very standardized.
        
    - **Performance:** Extremely low memory footprint, ideal for low-power devices.
        
- **Pros:** Easiest deployment, "Single Binary" architecture, modern/fast frontend.
    
- **Cons:** Go's error handling (`if err != nil`) is repetitive (though AI handles this well).
    
- **Deployment Details:**
    
    - **Docker:** Multi-stage build. Stage 1: Build Vue `dist`. Stage 2: `go build` with embedded dist. Stage 3: Scratch/Alpine container.
        
    - **CasaOS:** Simple `docker-compose.yml` mounting a `/config` volume for the `sqlite` database.
        
- **Data Store:** **SQLite** (via `modernc.org/sqlite` for a pure Go implementation without CGO, making cross-compilation trivial).

https://www.youtube.com/watch?v=PT8alMw3GFI

---

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- yt-dlp
- ffmpeg

### Quick Start with Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/listenbucket/listenbucket.git
cd listenbucket

# Build and run with Docker Compose
docker-compose up -d

# Access the app at http://localhost:8080
```

### Development Setup

```bash
# Install dependencies
make deps

# Run in development mode (backend on :8080, frontend on :5173)
make dev

# Or run separately:
# Terminal 1: Backend
make dev-backend

# Terminal 2: Frontend (with hot reload)
make dev-frontend
```

### Build from Source

```bash
# Build everything (frontend + backend)
make build

# Run
./bin/listenbucket
```

### Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DATA_DIR` | `./data` | Directory for database and media files |
| `BASE_URL` | `http://localhost:8080` | Public URL (used in RSS feeds) |

### Project Structure

```
.
├── cmd/listenbucket/      # Application entry point
├── internal/
│   ├── api/               # HTTP handlers and routing
│   ├── config/            # Configuration
│   ├── database/          # SQLite database layer
│   ├── downloader/        # yt-dlp integration
│   └── podcast/           # RSS feed generation
├── web/                   # Vue.js frontend
│   ├── src/
│   │   ├── components/    # Vue components
│   │   ├── views/         # Page views
│   │   ├── stores/        # Pinia stores
│   │   └── router/        # Vue Router
│   └── ...
├── Dockerfile             # Multi-stage Docker build
├── docker-compose.yml     # Docker Compose config
└── Makefile              # Build commands
```

### CasaOS Installation

Use the `docker-compose.casaos.yml` file or install directly from Docker Hub:

```bash
docker pull listenbucket/listenbucket:latest
```

### Running Tests

```bash
make test
```

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/feeds` | List all feeds |
| POST | `/api/feeds` | Create a new feed |
| GET | `/api/feeds/:id` | Get feed details |
| DELETE | `/api/feeds/:id` | Delete a feed |
| GET | `/api/feeds/:id/rss` | Get RSS feed XML |
| GET | `/api/feeds/:id/episodes` | List episodes |
| POST | `/api/feeds/:id/sources` | Add a source |
| GET | `/api/feeds/:id/sources` | List sources |
| DELETE | `/api/episodes/:id` | Delete episode |
| DELETE | `/api/sources/:id` | Delete source |
| GET | `/api/media/:filename` | Serve media files |

---

## License

MIT