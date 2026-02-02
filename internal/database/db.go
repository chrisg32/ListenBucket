package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

func New(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := conn.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	// Ensure default feed exists
	if err := db.ensureDefaultFeed(); err != nil {
		return nil, fmt.Errorf("failed to create default feed: %w", err)
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS feeds (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT DEFAULT '',
		image_url TEXT DEFAULT '',
		is_default INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS episodes (
		id TEXT PRIMARY KEY,
		feed_id TEXT NOT NULL,
		source_id TEXT DEFAULT '',
		title TEXT NOT NULL,
		description TEXT DEFAULT '',
		image_url TEXT DEFAULT '',
		audio_url TEXT DEFAULT '',
		duration INTEGER DEFAULT 0,
		source_url TEXT DEFAULT '',
		status TEXT DEFAULT 'pending',
		error_msg TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS sources (
		id TEXT PRIMARY KEY,
		feed_id TEXT NOT NULL,
		url TEXT NOT NULL,
		type TEXT DEFAULT 'video',
		title TEXT DEFAULT '',
		image_url TEXT DEFAULT '',
		last_checked DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_episodes_feed_id ON episodes(feed_id);
	CREATE INDEX IF NOT EXISTS idx_episodes_status ON episodes(status);
	CREATE INDEX IF NOT EXISTS idx_sources_feed_id ON sources(feed_id);
	`

	_, err := db.conn.Exec(schema)
	return err
}

func (db *DB) ensureDefaultFeed() error {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM feeds WHERE is_default = 1").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = db.conn.Exec(
			"INSERT INTO feeds (id, title, description, is_default, created_at, updated_at) VALUES (?, ?, ?, 1, ?, ?)",
			uuid.New().String(),
			"Listen Later",
			"Your personal listen later feed",
			time.Now(),
			time.Now(),
		)
	}
	return err
}

// Feed operations

func (db *DB) GetFeeds() ([]Feed, error) {
	rows, err := db.conn.Query("SELECT id, title, description, image_url, is_default, created_at, updated_at FROM feeds ORDER BY is_default DESC, created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []Feed
	for rows.Next() {
		var f Feed
		if err := rows.Scan(&f.ID, &f.Title, &f.Description, &f.ImageURL, &f.IsDefault, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		feeds = append(feeds, f)
	}
	return feeds, rows.Err()
}

func (db *DB) GetFeed(id string) (*Feed, error) {
	var f Feed
	err := db.conn.QueryRow(
		"SELECT id, title, description, image_url, is_default, created_at, updated_at FROM feeds WHERE id = ?",
		id,
	).Scan(&f.ID, &f.Title, &f.Description, &f.ImageURL, &f.IsDefault, &f.CreatedAt, &f.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (db *DB) CreateFeed(title, description string) (*Feed, error) {
	f := Feed{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		IsDefault:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := db.conn.Exec(
		"INSERT INTO feeds (id, title, description, is_default, created_at, updated_at) VALUES (?, ?, ?, 0, ?, ?)",
		f.ID, f.Title, f.Description, f.CreatedAt, f.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (db *DB) UpdateFeed(id, title, description string) error {
	_, err := db.conn.Exec(
		"UPDATE feeds SET title = ?, description = ?, updated_at = ? WHERE id = ? AND is_default = 0",
		title, description, time.Now(), id,
	)
	return err
}

func (db *DB) UpdateFeedImage(id, imageURL string) error {
	_, err := db.conn.Exec(
		"UPDATE feeds SET image_url = ?, updated_at = ? WHERE id = ?",
		imageURL, time.Now(), id,
	)
	return err
}

func (db *DB) DeleteFeed(id string) error {
	// Don't allow deleting the default feed
	_, err := db.conn.Exec("DELETE FROM feeds WHERE id = ? AND is_default = 0", id)
	return err
}

// Episode operations

func (db *DB) GetEpisodesByFeed(feedID string) ([]Episode, error) {
	rows, err := db.conn.Query(
		"SELECT id, feed_id, source_id, title, description, image_url, audio_url, duration, source_url, status, error_msg, created_at FROM episodes WHERE feed_id = ? ORDER BY created_at DESC",
		feedID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var episodes []Episode
	for rows.Next() {
		var e Episode
		if err := rows.Scan(&e.ID, &e.FeedID, &e.SourceID, &e.Title, &e.Description, &e.ImageURL, &e.AudioURL, &e.Duration, &e.SourceURL, &e.Status, &e.ErrorMsg, &e.CreatedAt); err != nil {
			return nil, err
		}
		episodes = append(episodes, e)
	}
	return episodes, rows.Err()
}

func (db *DB) GetEpisode(id string) (*Episode, error) {
	var e Episode
	err := db.conn.QueryRow(
		"SELECT id, feed_id, source_id, title, description, image_url, audio_url, duration, source_url, status, error_msg, created_at FROM episodes WHERE id = ?",
		id,
	).Scan(&e.ID, &e.FeedID, &e.SourceID, &e.Title, &e.Description, &e.ImageURL, &e.AudioURL, &e.Duration, &e.SourceURL, &e.Status, &e.ErrorMsg, &e.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (db *DB) CreateEpisode(feedID, sourceID, title, description, imageURL, sourceURL string) (*Episode, error) {
	e := Episode{
		ID:          uuid.New().String(),
		FeedID:      feedID,
		SourceID:    sourceID,
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
		SourceURL:   sourceURL,
		Status:      StatusPending,
		CreatedAt:   time.Now(),
	}

	_, err := db.conn.Exec(
		"INSERT INTO episodes (id, feed_id, source_id, title, description, image_url, source_url, status, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		e.ID, e.FeedID, e.SourceID, e.Title, e.Description, e.ImageURL, e.SourceURL, e.Status, e.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (db *DB) UpdateEpisodeStatus(id, status, errorMsg string) error {
	_, err := db.conn.Exec(
		"UPDATE episodes SET status = ?, error_msg = ? WHERE id = ?",
		status, errorMsg, id,
	)
	return err
}

func (db *DB) UpdateEpisodeAudio(id, audioURL string, duration int) error {
	_, err := db.conn.Exec(
		"UPDATE episodes SET audio_url = ?, duration = ?, status = ? WHERE id = ?",
		audioURL, duration, StatusReady, id,
	)
	return err
}

func (db *DB) DeleteEpisode(id string) error {
	_, err := db.conn.Exec("DELETE FROM episodes WHERE id = ?", id)
	return err
}

func (db *DB) GetPendingEpisodes() ([]Episode, error) {
	rows, err := db.conn.Query(
		"SELECT id, feed_id, source_id, title, description, image_url, audio_url, duration, source_url, status, error_msg, created_at FROM episodes WHERE status = ? LIMIT 10",
		StatusPending,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var episodes []Episode
	for rows.Next() {
		var e Episode
		if err := rows.Scan(&e.ID, &e.FeedID, &e.SourceID, &e.Title, &e.Description, &e.ImageURL, &e.AudioURL, &e.Duration, &e.SourceURL, &e.Status, &e.ErrorMsg, &e.CreatedAt); err != nil {
			return nil, err
		}
		episodes = append(episodes, e)
	}
	return episodes, rows.Err()
}

// Source operations

func (db *DB) GetSourcesByFeed(feedID string) ([]Source, error) {
	rows, err := db.conn.Query(
		"SELECT id, feed_id, url, type, title, image_url, last_checked, created_at FROM sources WHERE feed_id = ? ORDER BY created_at DESC",
		feedID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []Source
	for rows.Next() {
		var s Source
		var lastChecked sql.NullTime
		if err := rows.Scan(&s.ID, &s.FeedID, &s.URL, &s.Type, &s.Title, &s.ImageURL, &lastChecked, &s.CreatedAt); err != nil {
			return nil, err
		}
		if lastChecked.Valid {
			s.LastChecked = lastChecked.Time
		}
		sources = append(sources, s)
	}
	return sources, rows.Err()
}

func (db *DB) CreateSource(feedID, url, sourceType, title, imageURL string) (*Source, error) {
	s := Source{
		ID:        uuid.New().String(),
		FeedID:    feedID,
		URL:       url,
		Type:      sourceType,
		Title:     title,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
	}

	_, err := db.conn.Exec(
		"INSERT INTO sources (id, feed_id, url, type, title, image_url, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		s.ID, s.FeedID, s.URL, s.Type, s.Title, s.ImageURL, s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (db *DB) UpdateSourceLastChecked(id string) error {
	_, err := db.conn.Exec("UPDATE sources SET last_checked = ? WHERE id = ?", time.Now(), id)
	return err
}

func (db *DB) DeleteSource(id string) error {
	_, err := db.conn.Exec("DELETE FROM sources WHERE id = ?", id)
	return err
}

func (db *DB) GetSourcesNeedingUpdate() ([]Source, error) {
	// Get sources that haven't been checked in the last hour
	rows, err := db.conn.Query(
		"SELECT id, feed_id, url, type, title, image_url, last_checked, created_at FROM sources WHERE type IN (?, ?) AND (last_checked IS NULL OR last_checked < datetime('now', '-1 hour'))",
		SourceTypePlaylist, SourceTypeChannel,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []Source
	for rows.Next() {
		var s Source
		var lastChecked sql.NullTime
		if err := rows.Scan(&s.ID, &s.FeedID, &s.URL, &s.Type, &s.Title, &s.ImageURL, &lastChecked, &s.CreatedAt); err != nil {
			return nil, err
		}
		if lastChecked.Valid {
			s.LastChecked = lastChecked.Time
		}
		sources = append(sources, s)
	}
	return sources, rows.Err()
}

func (db *DB) EpisodeExistsForSource(sourceURL, feedID string) (bool, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM episodes WHERE source_url = ? AND feed_id = ?", sourceURL, feedID).Scan(&count)
	return count > 0, err
}

func (db *DB) GetReadyEpisodesByFeed(feedID string) ([]Episode, error) {
	rows, err := db.conn.Query(
		"SELECT id, feed_id, source_id, title, description, image_url, audio_url, duration, source_url, status, error_msg, created_at FROM episodes WHERE feed_id = ? AND status = ? ORDER BY created_at DESC",
		feedID, StatusReady,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var episodes []Episode
	for rows.Next() {
		var e Episode
		if err := rows.Scan(&e.ID, &e.FeedID, &e.SourceID, &e.Title, &e.Description, &e.ImageURL, &e.AudioURL, &e.Duration, &e.SourceURL, &e.Status, &e.ErrorMsg, &e.CreatedAt); err != nil {
			return nil, err
		}
		episodes = append(episodes, e)
	}
	return episodes, rows.Err()
}

// GetAllEpisodes returns all episodes across all feeds
func (db *DB) GetAllEpisodes() ([]Episode, error) {
	rows, err := db.conn.Query(
		"SELECT id, feed_id, source_id, title, description, image_url, audio_url, duration, source_url, status, error_msg, created_at FROM episodes ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var episodes []Episode
	for rows.Next() {
		var e Episode
		if err := rows.Scan(&e.ID, &e.FeedID, &e.SourceID, &e.Title, &e.Description, &e.ImageURL, &e.AudioURL, &e.Duration, &e.SourceURL, &e.Status, &e.ErrorMsg, &e.CreatedAt); err != nil {
			return nil, err
		}
		episodes = append(episodes, e)
	}
	return episodes, rows.Err()
}

// UpdateEpisodeMetadata updates the title and description of an episode
func (db *DB) UpdateEpisodeMetadata(id, title, description string) error {
	_, err := db.conn.Exec(
		"UPDATE episodes SET title = ?, description = ? WHERE id = ?",
		title, description, id,
	)
	return err
}

// GetAllSources returns all sources across all feeds
func (db *DB) GetAllSources() ([]Source, error) {
	rows, err := db.conn.Query(
		"SELECT id, feed_id, url, type, title, image_url, last_checked, created_at FROM sources ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []Source
	for rows.Next() {
		var s Source
		var lastChecked sql.NullTime
		if err := rows.Scan(&s.ID, &s.FeedID, &s.URL, &s.Type, &s.Title, &s.ImageURL, &lastChecked, &s.CreatedAt); err != nil {
			return nil, err
		}
		if lastChecked.Valid {
			s.LastChecked = lastChecked.Time
		}
		sources = append(sources, s)
	}
	return sources, rows.Err()
}

// GetSource returns a single source by ID
func (db *DB) GetSource(id string) (*Source, error) {
	var s Source
	var lastChecked sql.NullTime
	err := db.conn.QueryRow(
		"SELECT id, feed_id, url, type, title, image_url, last_checked, created_at FROM sources WHERE id = ?",
		id,
	).Scan(&s.ID, &s.FeedID, &s.URL, &s.Type, &s.Title, &s.ImageURL, &lastChecked, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if lastChecked.Valid {
		s.LastChecked = lastChecked.Time
	}
	return &s, nil
}
