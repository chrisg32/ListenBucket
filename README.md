# ListenBucket

[![CI](https://github.com/listenbucket/listenbucket/actions/workflows/ci.yml/badge.svg)](https://github.com/listenbucket/listenbucket/actions/workflows/ci.yml)
[![Release](https://github.com/listenbucket/listenbucket/actions/workflows/release.yml/badge.svg)](https://github.com/listenbucket/listenbucket/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/listenbucket/listenbucket)](https://goreportcard.com/report/github.com/listenbucket/listenbucket)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Turn YouTube videos into podcast feeds.** ListenBucket lets you add YouTube videos, playlists, and channels to personal podcast feeds that you can subscribe to in any podcast app.

![ListenBucket Screenshot](listen-later.png)

## Features

- **YouTube Support** - Add individual videos, entire playlists, or channels
- **Automatic Updates** - Playlists and channels are checked hourly for new content
- **Podcast Compatible** - Subscribe in Apple Podcasts, Overcast, Pocket Casts, or any podcast app
- **Self-Hosted** - Your data stays on your server
- **Single Binary** - Easy deployment with Docker or standalone
- **Dark Mode** - Light and dark themes
- **Mobile Friendly** - Responsive design works on any device

## Quick Start

### Docker (Recommended)

```bash
docker run -d \
  --name listenbucket \
  -p 8080:8080 \
  -v listenbucket_data:/data \
  listenbucket/listenbucket:latest
```

Or with Docker Compose:

```bash
curl -O https://raw.githubusercontent.com/listenbucket/listenbucket/main/docker-compose.yml
docker-compose up -d
```

Access the app at **http://localhost:8080**

### CasaOS

Install directly from Docker Hub or use the included `docker-compose.casaos.yml`.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DATA_DIR` | `/data` | Database and media storage |
| `BASE_URL` | `http://localhost:8080` | Public URL for RSS feeds |

## Usage

1. **Create a Feed** - Or use the default "Listen Later" feed
2. **Add Sources** - Paste YouTube video, playlist, or channel URLs
3. **Subscribe** - Copy the RSS link or open directly in your podcast app
4. **Listen** - Episodes download automatically as MP3s

## Development

### Prerequisites

- Go 1.22+
- Node.js 20+
- yt-dlp
- ffmpeg

### Setup

```bash
# Clone the repository
git clone https://github.com/listenbucket/listenbucket.git
cd listenbucket

# Install dependencies
make deps

# Run in development mode
make dev
```

### Build

```bash
# Build everything
make build

# Run tests
make test

# Build Docker image
docker build -t listenbucket .
```

## API

ListenBucket provides a RESTful API for automation and integration.

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/feeds` | List all feeds |
| `POST` | `/api/feeds` | Create a feed |
| `GET` | `/api/feeds/{id}/rss` | Get RSS feed |
| `POST` | `/api/feeds/{id}/sources` | Add a source |
| `GET` | `/api/episodes` | List all episodes |

See the full [API documentation](#api-reference) below.

<details>
<summary><strong>API Reference</strong></summary>

### Feeds

```bash
# List feeds
curl http://localhost:8080/api/feeds

# Create feed
curl -X POST http://localhost:8080/api/feeds \
  -H "Content-Type: application/json" \
  -d '{"title":"My Feed"}'

# Get RSS
curl http://localhost:8080/api/feeds/{feedId}/rss
```

### Sources

```bash
# Add YouTube video
curl -X POST http://localhost:8080/api/feeds/{feedId}/sources \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.youtube.com/watch?v=..."}'

# Add playlist (new videos only)
curl -X POST http://localhost:8080/api/feeds/{feedId}/sources \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.youtube.com/playlist?list=...","include_back_catalog":false}'
```

### Episodes

```bash
# List episodes
curl http://localhost:8080/api/feeds/{feedId}/episodes

# Retry failed download
curl -X POST http://localhost:8080/api/episodes/{episodeId}/retry

# Delete episode
curl -X DELETE http://localhost:8080/api/episodes/{episodeId}
```

### Error Responses

```json
{
  "error": "not_found",
  "message": "feed not found",
  "code": 404
}
```

</details>

## Architecture

ListenBucket uses a single-binary architecture with:
- **Backend**: Go with Chi router and SQLite
- **Frontend**: Vue.js 3 with Tailwind CSS (embedded in binary)
- **Media**: yt-dlp for downloading, ffmpeg for audio extraction

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed technical documentation.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the discontinued ListenBox app
- Architecture inspired by [Mailpit](https://mailpit.axllent.org/)
- Powered by [yt-dlp](https://github.com/yt-dlp/yt-dlp)
