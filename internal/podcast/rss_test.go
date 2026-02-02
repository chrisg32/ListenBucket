package podcast

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"github.com/listenbucket/listenbucket/internal/database"
)

func TestGenerateFeed(t *testing.T) {
	feed := &database.Feed{
		ID:          "test-feed-id",
		Title:       "Test Feed",
		Description: "Test Description",
		ImageURL:    "http://example.com/image.jpg",
		IsDefault:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	episodes := []database.Episode{
		{
			ID:          "ep-1",
			FeedID:      feed.ID,
			Title:       "Episode 1",
			Description: "First episode",
			ImageURL:    "http://example.com/ep1.jpg",
			AudioURL:    "http://example.com/ep1.mp3",
			Duration:    3600, // 1 hour
			SourceURL:   "http://youtube.com/watch?v=123",
			Status:      database.StatusReady,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "ep-2",
			FeedID:      feed.ID,
			Title:       "Episode 2",
			Description: "Second episode",
			ImageURL:    "http://example.com/ep2.jpg",
			AudioURL:    "http://example.com/ep2.mp3",
			Duration:    1800, // 30 minutes
			SourceURL:   "http://youtube.com/watch?v=456",
			Status:      database.StatusReady,
			CreatedAt:   time.Now(),
		},
	}

	baseURL := "http://localhost:8080"

	rss, err := GenerateFeed(feed, episodes, baseURL)
	if err != nil {
		t.Fatalf("failed to generate feed: %v", err)
	}

	// Verify it's valid XML
	var result RSS
	if err := xml.Unmarshal(rss, &result); err != nil {
		t.Fatalf("generated feed is not valid XML: %v", err)
	}

	// Verify basic structure
	if result.Channel.Title != feed.Title {
		t.Errorf("expected title '%s', got '%s'", feed.Title, result.Channel.Title)
	}

	if result.Channel.Description != feed.Description {
		t.Errorf("expected description '%s', got '%s'", feed.Description, result.Channel.Description)
	}

	if len(result.Channel.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Channel.Items))
	}

	// Verify first episode
	item := result.Channel.Items[0]
	if item.Title != "Episode 1" {
		t.Errorf("expected item title 'Episode 1', got '%s'", item.Title)
	}

	if item.Enclosure == nil {
		t.Error("expected enclosure to be set")
	} else if item.Enclosure.URL != "http://example.com/ep1.mp3" {
		t.Errorf("expected enclosure URL 'http://example.com/ep1.mp3', got '%s'", item.Enclosure.URL)
	}

	// Verify duration formatting (1:00:00)
	if item.ITunesDuration != "1:00:00" {
		t.Errorf("expected duration '1:00:00', got '%s'", item.ITunesDuration)
	}

	// Verify XML header
	if !strings.HasPrefix(string(rss), "<?xml") {
		t.Error("expected XML header")
	}
}

func TestGenerateEmptyFeed(t *testing.T) {
	feed := &database.Feed{
		ID:          "empty-feed",
		Title:       "Empty Feed",
		Description: "",
		UpdatedAt:   time.Now(),
	}

	rss, err := GenerateFeed(feed, []database.Episode{}, "http://localhost:8080")
	if err != nil {
		t.Fatalf("failed to generate empty feed: %v", err)
	}

	var result RSS
	if err := xml.Unmarshal(rss, &result); err != nil {
		t.Fatalf("generated feed is not valid XML: %v", err)
	}

	if len(result.Channel.Items) != 0 {
		t.Errorf("expected 0 items for empty feed, got %d", len(result.Channel.Items))
	}
}

func TestDurationFormatting(t *testing.T) {
	feed := &database.Feed{
		ID:        "test",
		Title:     "Test",
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		duration int
		expected string
	}{
		{60, "1:00"},           // 1 minute
		{90, "1:30"},           // 1:30
		{3600, "1:00:00"},      // 1 hour
		{3665, "1:01:05"},      // 1:01:05
		{7384, "2:03:04"},      // 2:03:04
	}

	for _, tt := range tests {
		episodes := []database.Episode{
			{
				ID:       "test",
				Duration: tt.duration,
				AudioURL: "http://example.com/audio.mp3",
			},
		}

		rss, _ := GenerateFeed(feed, episodes, "http://localhost:8080")

		var result RSS
		xml.Unmarshal(rss, &result)

		if len(result.Channel.Items) == 0 {
			t.Fatal("no items in feed")
		}

		if result.Channel.Items[0].ITunesDuration != tt.expected {
			t.Errorf("duration %d: expected '%s', got '%s'", tt.duration, tt.expected, result.Channel.Items[0].ITunesDuration)
		}
	}
}
