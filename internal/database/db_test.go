package database

import (
	"os"
	"testing"
)

func setupTestDB(t *testing.T) (*DB, func()) {
	t.Helper()

	// Create temp file for test database
	tmpFile, err := os.CreateTemp("", "listenbucket-test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()

	db, err := New(tmpFile.Name())
	if err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to create database: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(tmpFile.Name())
	}

	return db, cleanup
}

func TestNewDatabase(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	if db == nil {
		t.Fatal("expected non-nil database")
	}
}

func TestDefaultFeedCreation(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}

	if len(feeds) != 1 {
		t.Fatalf("expected 1 default feed, got %d", len(feeds))
	}

	if feeds[0].Title != "Listen Later" {
		t.Errorf("expected title 'Listen Later', got '%s'", feeds[0].Title)
	}

	if !feeds[0].IsDefault {
		t.Error("expected default feed to have IsDefault=true")
	}
}

func TestCreateFeed(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	feed, err := db.CreateFeed("Test Feed", "Test Description")
	if err != nil {
		t.Fatalf("failed to create feed: %v", err)
	}

	if feed.Title != "Test Feed" {
		t.Errorf("expected title 'Test Feed', got '%s'", feed.Title)
	}

	if feed.Description != "Test Description" {
		t.Errorf("expected description 'Test Description', got '%s'", feed.Description)
	}

	if feed.IsDefault {
		t.Error("expected new feed to not be default")
	}

	// Verify feed was persisted
	retrieved, err := db.GetFeed(feed.ID)
	if err != nil {
		t.Fatalf("failed to get feed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected to retrieve created feed")
	}

	if retrieved.Title != "Test Feed" {
		t.Errorf("expected retrieved title 'Test Feed', got '%s'", retrieved.Title)
	}
}

func TestDeleteFeed(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a non-default feed
	feed, err := db.CreateFeed("Deletable Feed", "")
	if err != nil {
		t.Fatalf("failed to create feed: %v", err)
	}

	// Delete it
	err = db.DeleteFeed(feed.ID)
	if err != nil {
		t.Fatalf("failed to delete feed: %v", err)
	}

	// Verify it's gone
	retrieved, err := db.GetFeed(feed.ID)
	if err != nil {
		t.Fatalf("failed to get feed: %v", err)
	}

	if retrieved != nil {
		t.Error("expected feed to be deleted")
	}
}

func TestCannotDeleteDefaultFeed(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}

	defaultFeed := feeds[0]
	if !defaultFeed.IsDefault {
		t.Fatal("expected first feed to be default")
	}

	// Try to delete default feed
	err = db.DeleteFeed(defaultFeed.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it still exists
	retrieved, err := db.GetFeed(defaultFeed.ID)
	if err != nil {
		t.Fatalf("failed to get feed: %v", err)
	}

	if retrieved == nil {
		t.Error("default feed should not be deletable")
	}
}

func TestCreateEpisode(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	feeds, _ := db.GetFeeds()
	feedID := feeds[0].ID

	episode, err := db.CreateEpisode(feedID, "", "Test Episode", "Description", "http://example.com/thumb.jpg", "http://youtube.com/watch?v=123")
	if err != nil {
		t.Fatalf("failed to create episode: %v", err)
	}

	if episode.Title != "Test Episode" {
		t.Errorf("expected title 'Test Episode', got '%s'", episode.Title)
	}

	if episode.Status != StatusPending {
		t.Errorf("expected status '%s', got '%s'", StatusPending, episode.Status)
	}

	// Verify episode was persisted
	retrieved, err := db.GetEpisode(episode.ID)
	if err != nil {
		t.Fatalf("failed to get episode: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected to retrieve created episode")
	}
}

func TestUpdateEpisodeStatus(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	feeds, _ := db.GetFeeds()
	feedID := feeds[0].ID

	episode, _ := db.CreateEpisode(feedID, "", "Test Episode", "", "", "http://youtube.com/watch?v=123")

	// Update to downloading
	err := db.UpdateEpisodeStatus(episode.ID, StatusDownloading, "")
	if err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	retrieved, _ := db.GetEpisode(episode.ID)
	if retrieved.Status != StatusDownloading {
		t.Errorf("expected status '%s', got '%s'", StatusDownloading, retrieved.Status)
	}

	// Update to error
	err = db.UpdateEpisodeStatus(episode.ID, StatusError, "Download failed")
	if err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	retrieved, _ = db.GetEpisode(episode.ID)
	if retrieved.Status != StatusError {
		t.Errorf("expected status '%s', got '%s'", StatusError, retrieved.Status)
	}
	if retrieved.ErrorMsg != "Download failed" {
		t.Errorf("expected error message 'Download failed', got '%s'", retrieved.ErrorMsg)
	}
}

func TestGetPendingEpisodes(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	feeds, _ := db.GetFeeds()
	feedID := feeds[0].ID

	// Create multiple episodes
	_, _ = db.CreateEpisode(feedID, "", "Pending 1", "", "", "http://youtube.com/watch?v=1")
	_, _ = db.CreateEpisode(feedID, "", "Pending 2", "", "", "http://youtube.com/watch?v=2")

	ep3, _ := db.CreateEpisode(feedID, "", "Not Pending", "", "", "http://youtube.com/watch?v=3")
	_ = db.UpdateEpisodeStatus(ep3.ID, StatusReady, "")

	pending, err := db.GetPendingEpisodes()
	if err != nil {
		t.Fatalf("failed to get pending episodes: %v", err)
	}

	if len(pending) != 2 {
		t.Errorf("expected 2 pending episodes, got %d", len(pending))
	}
}

func TestCreateSource(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	feeds, _ := db.GetFeeds()
	feedID := feeds[0].ID

	source, err := db.CreateSource(feedID, "http://youtube.com/playlist?list=123", SourceTypePlaylist, "Test Playlist", "", true)
	if err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	if source.Type != SourceTypePlaylist {
		t.Errorf("expected type '%s', got '%s'", SourceTypePlaylist, source.Type)
	}

	if !source.IncludeBackCatalog {
		t.Error("expected IncludeBackCatalog to be true")
	}

	// Verify source was persisted
	retrieved, err := db.GetSource(source.ID)
	if err != nil {
		t.Fatalf("failed to get source: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected to retrieve created source")
	}
}

func TestEpisodeExistsForSource(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	feeds, _ := db.GetFeeds()
	feedID := feeds[0].ID

	sourceURL := "http://youtube.com/watch?v=unique123"

	// Should not exist initially
	exists, err := db.EpisodeExistsForSource(sourceURL, feedID)
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if exists {
		t.Error("expected episode to not exist")
	}

	// Create episode
	_, _ = db.CreateEpisode(feedID, "", "Test", "", "", sourceURL)

	// Should exist now
	exists, err = db.EpisodeExistsForSource(sourceURL, feedID)
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if !exists {
		t.Error("expected episode to exist")
	}
}

func TestGetReadyEpisodesByFeed(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	feeds, _ := db.GetFeeds()
	feedID := feeds[0].ID

	// Create episodes with different statuses
	ep1, _ := db.CreateEpisode(feedID, "", "Ready Episode", "", "", "http://youtube.com/watch?v=1")
	_ = db.UpdateEpisodeAudio(ep1.ID, "http://example.com/audio.mp3", 300)

	_, _ = db.CreateEpisode(feedID, "", "Pending Episode", "", "", "http://youtube.com/watch?v=2")

	ready, err := db.GetReadyEpisodesByFeed(feedID)
	if err != nil {
		t.Fatalf("failed to get ready episodes: %v", err)
	}

	if len(ready) != 1 {
		t.Errorf("expected 1 ready episode, got %d", len(ready))
	}

	if ready[0].Title != "Ready Episode" {
		t.Errorf("expected 'Ready Episode', got '%s'", ready[0].Title)
	}
}
