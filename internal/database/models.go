package database

import "time"

type Feed struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Episode struct {
	ID          string    `json:"id"`
	FeedID      string    `json:"feed_id"`
	SourceID    string    `json:"source_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	AudioURL    string    `json:"audio_url"`
	Duration    int       `json:"duration"` // seconds
	SourceURL   string    `json:"source_url"`
	Status      string    `json:"status"` // pending, downloading, ready, error
	ErrorMsg    string    `json:"error_msg"`
	CreatedAt   time.Time `json:"created_at"`
}

type Source struct {
	ID                 string    `json:"id"`
	FeedID             string    `json:"feed_id"`
	URL                string    `json:"url"`
	Type               string    `json:"type"` // video, playlist, channel
	Title              string    `json:"title"`
	ImageURL           string    `json:"image_url"`
	IncludeBackCatalog bool      `json:"include_back_catalog"` // For playlists/channels: add existing videos
	LastChecked        time.Time `json:"last_checked"`
	CreatedAt          time.Time `json:"created_at"`
}

const (
	StatusPending     = "pending"
	StatusDownloading = "downloading"
	StatusReady       = "ready"
	StatusError       = "error"
)

const (
	SourceTypeVideo    = "video"
	SourceTypePlaylist = "playlist"
	SourceTypeChannel  = "channel"
)
