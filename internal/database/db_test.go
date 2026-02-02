package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "listenbucket-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	// Verify default feed was created
	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}

	if len(feeds) != 1 {
		t.Errorf("expected 1 default feed, got %d", len(feeds))
	}

	if feeds[0].Title != "Listen Later" {
		t.Errorf("expected default feed title 'Listen Later', got '%s'", feeds[0].Title)
	}

	if !feeds[0].IsDefault {
		t.Error("expected default feed to be marked as default")
	}
}

func TestCreateAndDeleteFeed(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "listenbucket-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	db, err := New(filepath.Join(tmpDir, "test.db"))
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	// Create a new feed
	feed, err := db.CreateFeed("Test Feed", "Test Description")
	if err != nil {
		t.Fatalf("failed to create feed: %v", err)
	}

	if feed.Title != "Test Feed" {
		t.Errorf("expected title 'Test Feed', got '%s'", feed.Title)
	}

	if feed.IsDefault {
		t.Error("new feed should not be default")
	}

	// Verify we now have 2 feeds
	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}

	if len(feeds) != 2 {
		t.Errorf("expected 2 feeds, got %d", len(feeds))
	}

	// Delete the feed
	err = db.DeleteFeed(feed.ID)
	if err != nil {
		t.Fatalf("failed to delete feed: %v", err)
	}

	// Verify we're back to 1 feed
	feeds, err = db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}

	if len(feeds) != 1 {
		t.Errorf("expected 1 feed after deletion, got %d", len(feeds))
	}
}

func TestCannotDeleteDefaultFeed(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "listenbucket-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	db, err := New(filepath.Join(tmpDir, "test.db"))
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	feeds, err := db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}

	// Try to delete the default feed
	err = db.DeleteFeed(feeds[0].ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Feed should still exist (DELETE WHERE is_default = 0 won't match)
	feeds, err = db.GetFeeds()
	if err != nil {
		t.Fatalf("failed to get feeds: %v", err)
	}

	if len(feeds) != 1 {
		t.Errorf("default feed should not be deletable, expected 1 feed, got %d", len(feeds))
	}
}

func TestEpisodeOperations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "listenbucket-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	db, err := New(filepath.Join(tmpDir, "test.db"))
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	feeds, _ := db.GetFeeds()
	feedID := feeds[0].ID

	// Create episode
	ep, err := db.CreateEpisode(feedID, "", "Test Episode", "Description", "http://example.com/thumb.jpg", "http://youtube.com/watch?v=123")
	if err != nil {
		t.Fatalf("failed to create episode: %v", err)
	}

	if ep.Status != StatusPending {
		t.Errorf("expected status %s, got %s", StatusPending, ep.Status)
	}

	// Get episodes
	episodes, err := db.GetEpisodesByFeed(feedID)
	if err != nil {
		t.Fatalf("failed to get episodes: %v", err)
	}

	if len(episodes) != 1 {
		t.Errorf("expected 1 episode, got %d", len(episodes))
	}

	// Update status
	err = db.UpdateEpisodeStatus(ep.ID, StatusDownloading, "")
	if err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	// Verify status updated
	updated, err := db.GetEpisode(ep.ID)
	if err != nil {
		t.Fatalf("failed to get episode: %v", err)
	}

	if updated.Status != StatusDownloading {
		t.Errorf("expected status %s, got %s", StatusDownloading, updated.Status)
	}

	// Update audio
	err = db.UpdateEpisodeAudio(ep.ID, "http://example.com/audio.mp3", 300)
	if err != nil {
		t.Fatalf("failed to update audio: %v", err)
	}

	// Verify ready episodes
	ready, err := db.GetReadyEpisodesByFeed(feedID)
	if err != nil {
		t.Fatalf("failed to get ready episodes: %v", err)
	}

	if len(ready) != 1 {
		t.Errorf("expected 1 ready episode, got %d", len(ready))
	}

	if ready[0].Duration != 300 {
		t.Errorf("expected duration 300, got %d", ready[0].Duration)
	}
}

func TestSourceOperations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "listenbucket-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	db, err := New(filepath.Join(tmpDir, "test.db"))
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	feeds, _ := db.GetFeeds()
	feedID := feeds[0].ID

	// Create source
	src, err := db.CreateSource(feedID, "http://youtube.com/watch?v=123", SourceTypeVideo, "Test Video", "http://example.com/thumb.jpg")
	if err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	if src.Type != SourceTypeVideo {
		t.Errorf("expected type %s, got %s", SourceTypeVideo, src.Type)
	}

	// Get sources
	sources, err := db.GetSourcesByFeed(feedID)
	if err != nil {
		t.Fatalf("failed to get sources: %v", err)
	}

	if len(sources) != 1 {
		t.Errorf("expected 1 source, got %d", len(sources))
	}

	// Delete source
	err = db.DeleteSource(src.ID)
	if err != nil {
		t.Fatalf("failed to delete source: %v", err)
	}

	sources, _ = db.GetSourcesByFeed(feedID)
	if len(sources) != 0 {
		t.Errorf("expected 0 sources after deletion, got %d", len(sources))
	}
}
