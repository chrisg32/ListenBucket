package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/listenbucket/listenbucket/internal/database"
)

type Downloader struct {
	db       *database.DB
	mediaDir string
	baseURL  string
}

type VideoInfo struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Thumbnail   string  `json:"thumbnail"`
	Duration    float64 `json:"duration"`
	WebpageURL  string  `json:"webpage_url"`
	Uploader    string  `json:"uploader"`
	PlaylistID  string  `json:"playlist_id"`
}

type PlaylistInfo struct {
	ID      string      `json:"id"`
	Title   string      `json:"title"`
	Entries []VideoInfo `json:"entries"`
}

func New(db *database.DB, mediaDir, baseURL string) *Downloader {
	return &Downloader{
		db:       db,
		mediaDir: mediaDir,
		baseURL:  baseURL,
	}
}

func (d *Downloader) Start(ctx context.Context) {
	// Process pending downloads every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Process immediately on start
	d.processPending(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.processPending(ctx)
		}
	}
}

func (d *Downloader) processPending(ctx context.Context) {
	episodes, err := d.db.GetPendingEpisodes()
	if err != nil {
		log.Printf("Error getting pending episodes: %v", err)
		return
	}

	for _, ep := range episodes {
		select {
		case <-ctx.Done():
			return
		default:
			if err := d.downloadEpisode(ctx, &ep); err != nil {
				log.Printf("Error downloading episode %s: %v", ep.ID, err)
				d.db.UpdateEpisodeStatus(ep.ID, database.StatusError, err.Error())
			}
		}
	}
}

func (d *Downloader) downloadEpisode(ctx context.Context, ep *database.Episode) error {
	d.db.UpdateEpisodeStatus(ep.ID, database.StatusDownloading, "")

	// Ensure media directory exists
	if err := os.MkdirAll(d.mediaDir, 0755); err != nil {
		return fmt.Errorf("failed to create media dir: %w", err)
	}

	outputPath := filepath.Join(d.mediaDir, ep.ID+".mp3")

	// Use yt-dlp to download and convert to mp3
	cmd := exec.CommandContext(ctx, "yt-dlp",
		"-x",                     // Extract audio
		"--audio-format", "mp3",  // Convert to mp3
		"--audio-quality", "0",   // Best quality
		"-o", outputPath,         // Output path
		"--no-playlist",          // Don't download playlist
		ep.SourceURL,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("yt-dlp failed: %w, output: %s", err, string(output))
	}

	// Get duration using ffprobe
	duration := d.getDuration(outputPath)

	// Set audio URL
	audioURL := fmt.Sprintf("%s/api/media/%s.mp3", d.baseURL, ep.ID)

	return d.db.UpdateEpisodeAudio(ep.ID, audioURL, duration)
}

func (d *Downloader) getDuration(path string) int {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	var duration float64
	fmt.Sscanf(strings.TrimSpace(string(output)), "%f", &duration)
	return int(duration)
}

func (d *Downloader) GetVideoInfo(url string) (*VideoInfo, error) {
	cmd := exec.Command("yt-dlp",
		"--dump-json",
		"--no-playlist",
		url,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp failed: %w", err)
	}

	var info VideoInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return &info, nil
}

func (d *Downloader) GetPlaylistInfo(url string) (*PlaylistInfo, []VideoInfo, error) {
	cmd := exec.Command("yt-dlp",
		"--dump-json",
		"--flat-playlist",
		url,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, nil, fmt.Errorf("yt-dlp failed: %w", err)
	}

	// yt-dlp outputs one JSON object per line for playlists
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var videos []VideoInfo

	for _, line := range lines {
		if line == "" {
			continue
		}
		var info VideoInfo
		if err := json.Unmarshal([]byte(line), &info); err != nil {
			continue
		}
		// For flat playlists, we need to construct the full URL
		if info.WebpageURL == "" && info.ID != "" {
			info.WebpageURL = "https://www.youtube.com/watch?v=" + info.ID
		}
		videos = append(videos, info)
	}

	playlist := &PlaylistInfo{
		Entries: videos,
	}

	// Try to get playlist metadata
	if len(videos) > 0 {
		// Extract playlist ID from URL
		playlist.ID = extractPlaylistID(url)
	}

	return playlist, videos, nil
}

func extractPlaylistID(url string) string {
	re := regexp.MustCompile(`[?&]list=([^&]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func (d *Downloader) DetectSourceType(url string) string {
	// Check for playlist
	if strings.Contains(url, "list=") {
		return database.SourceTypePlaylist
	}
	// Check for channel
	if strings.Contains(url, "/channel/") || strings.Contains(url, "/@") || strings.Contains(url, "/c/") || strings.Contains(url, "/user/") {
		return database.SourceTypeChannel
	}
	return database.SourceTypeVideo
}

func (d *Downloader) StartSourceChecker(ctx context.Context) {
	// Check sources every hour
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Check immediately on start
	d.checkSources(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.checkSources(ctx)
		}
	}
}

func (d *Downloader) checkSources(ctx context.Context) {
	sources, err := d.db.GetSourcesNeedingUpdate()
	if err != nil {
		log.Printf("Error getting sources: %v", err)
		return
	}

	for _, src := range sources {
		select {
		case <-ctx.Done():
			return
		default:
			d.checkSource(ctx, &src)
		}
	}
}

func (d *Downloader) checkSource(ctx context.Context, src *database.Source) {
	log.Printf("Checking source: %s (%s)", src.Title, src.URL)

	_, videos, err := d.GetPlaylistInfo(src.URL)
	if err != nil {
		log.Printf("Error checking source %s: %v", src.ID, err)
		return
	}

	for _, video := range videos {
		// Check if episode already exists
		exists, err := d.db.EpisodeExistsForSource(video.WebpageURL, src.FeedID)
		if err != nil {
			log.Printf("Error checking episode existence: %v", err)
			continue
		}
		if exists {
			continue
		}

		// Create new episode
		_, err = d.db.CreateEpisode(
			src.FeedID,
			src.ID,
			video.Title,
			video.Description,
			video.Thumbnail,
			video.WebpageURL,
		)
		if err != nil {
			log.Printf("Error creating episode: %v", err)
		}
	}

	d.db.UpdateSourceLastChecked(src.ID)
}

func (d *Downloader) AddSource(feedID, url string) (*database.Source, error) {
	sourceType := d.DetectSourceType(url)

	var title, imageURL string

	switch sourceType {
	case database.SourceTypeVideo:
		info, err := d.GetVideoInfo(url)
		if err != nil {
			return nil, err
		}
		title = info.Title
		imageURL = info.Thumbnail

		// Create source
		src, err := d.db.CreateSource(feedID, url, sourceType, title, imageURL)
		if err != nil {
			return nil, err
		}

		// Create episode directly for single video
		_, err = d.db.CreateEpisode(feedID, src.ID, info.Title, info.Description, info.Thumbnail, info.WebpageURL)
		if err != nil {
			return nil, err
		}

		// Update feed image if not set
		d.updateFeedImageIfNeeded(feedID, imageURL)

		return src, nil

	case database.SourceTypePlaylist, database.SourceTypeChannel:
		playlist, videos, err := d.GetPlaylistInfo(url)
		if err != nil {
			return nil, err
		}

		if len(videos) > 0 {
			title = fmt.Sprintf("Playlist: %s", playlist.ID)
			if videos[0].Thumbnail != "" {
				imageURL = videos[0].Thumbnail
			}
		}

		src, err := d.db.CreateSource(feedID, url, sourceType, title, imageURL)
		if err != nil {
			return nil, err
		}

		// Create episodes for all videos in playlist
		for _, video := range videos {
			_, err := d.db.CreateEpisode(feedID, src.ID, video.Title, video.Description, video.Thumbnail, video.WebpageURL)
			if err != nil {
				log.Printf("Error creating episode for video %s: %v", video.ID, err)
			}
		}

		d.updateFeedImageIfNeeded(feedID, imageURL)

		return src, nil
	}

	return nil, fmt.Errorf("unknown source type")
}

func (d *Downloader) updateFeedImageIfNeeded(feedID, imageURL string) {
	if imageURL == "" {
		return
	}

	feed, err := d.db.GetFeed(feedID)
	if err != nil || feed == nil {
		return
	}

	if feed.ImageURL == "" {
		d.db.UpdateFeedImage(feedID, imageURL)
	}
}

// RefreshSource manually triggers a check for new content from a source
func (d *Downloader) RefreshSource(src *database.Source) {
	log.Printf("Manually refreshing source: %s (%s)", src.Title, src.URL)

	if src.Type == database.SourceTypeVideo {
		// Single videos don't need refreshing
		d.db.UpdateSourceLastChecked(src.ID)
		return
	}

	_, videos, err := d.GetPlaylistInfo(src.URL)
	if err != nil {
		log.Printf("Error refreshing source %s: %v", src.ID, err)
		return
	}

	newCount := 0
	for _, video := range videos {
		exists, err := d.db.EpisodeExistsForSource(video.WebpageURL, src.FeedID)
		if err != nil {
			log.Printf("Error checking episode existence: %v", err)
			continue
		}
		if exists {
			continue
		}

		_, err = d.db.CreateEpisode(
			src.FeedID,
			src.ID,
			video.Title,
			video.Description,
			video.Thumbnail,
			video.WebpageURL,
		)
		if err != nil {
			log.Printf("Error creating episode: %v", err)
		} else {
			newCount++
		}
	}

	log.Printf("Refresh complete for source %s: %d new episodes", src.ID, newCount)
	d.db.UpdateSourceLastChecked(src.ID)
}
