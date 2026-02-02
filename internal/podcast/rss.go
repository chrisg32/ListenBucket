package podcast

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/listenbucket/listenbucket/internal/database"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	ITunes  string   `xml:"xmlns:itunes,attr"`
	Atom    string   `xml:"xmlns:atom,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string      `xml:"title"`
	Link        string      `xml:"link"`
	Description string      `xml:"description"`
	Language    string      `xml:"language"`
	Copyright   string      `xml:"copyright"`
	LastBuild   string      `xml:"lastBuildDate"`
	Image       *Image      `xml:"image,omitempty"`
	ITunesImage *ITunesImage `xml:"itunes:image,omitempty"`
	ITunesAuthor string     `xml:"itunes:author"`
	ITunesExplicit string   `xml:"itunes:explicit"`
	ITunesCategory *ITunesCategory `xml:"itunes:category,omitempty"`
	AtomLink    *AtomLink   `xml:"atom:link,omitempty"`
	Items       []Item      `xml:"item"`
}

type Image struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type ITunesImage struct {
	Href string `xml:"href,attr"`
}

type ITunesCategory struct {
	Text string `xml:"text,attr"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Item struct {
	Title       string       `xml:"title"`
	Link        string       `xml:"link,omitempty"`
	Description string       `xml:"description"`
	GUID        GUID         `xml:"guid"`
	PubDate     string       `xml:"pubDate"`
	Enclosure   *Enclosure   `xml:"enclosure,omitempty"`
	ITunesImage *ITunesImage `xml:"itunes:image,omitempty"`
	ITunesDuration string    `xml:"itunes:duration,omitempty"`
	ITunesExplicit string    `xml:"itunes:explicit"`
}

type GUID struct {
	IsPermaLink string `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

func GenerateFeed(feed *database.Feed, episodes []database.Episode, baseURL string) ([]byte, error) {
	feedURL := fmt.Sprintf("%s/api/feeds/%s/rss", baseURL, feed.ID)

	rss := RSS{
		Version: "2.0",
		ITunes:  "http://www.itunes.com/dtds/podcast-1.0.dtd",
		Atom:    "http://www.w3.org/2005/Atom",
		Channel: Channel{
			Title:       feed.Title,
			Link:        feedURL,
			Description: feed.Description,
			Language:    "en-us",
			Copyright:   fmt.Sprintf("© %d", time.Now().Year()),
			LastBuild:   feed.UpdatedAt.Format(time.RFC1123Z),
			ITunesAuthor: "ListenBucket",
			ITunesExplicit: "no",
			ITunesCategory: &ITunesCategory{Text: "Technology"},
			AtomLink: &AtomLink{
				Href: feedURL,
				Rel:  "self",
				Type: "application/rss+xml",
			},
		},
	}

	if feed.ImageURL != "" {
		rss.Channel.Image = &Image{
			URL:   feed.ImageURL,
			Title: feed.Title,
			Link:  feedURL,
		}
		rss.Channel.ITunesImage = &ITunesImage{Href: feed.ImageURL}
	}

	for _, ep := range episodes {
		item := Item{
			Title:       ep.Title,
			Link:        ep.SourceURL,
			Description: ep.Description,
			GUID: GUID{
				IsPermaLink: "false",
				Value:       ep.ID,
			},
			PubDate:        ep.CreatedAt.Format(time.RFC1123Z),
			ITunesExplicit: "no",
		}

		if ep.AudioURL != "" {
			item.Enclosure = &Enclosure{
				URL:    ep.AudioURL,
				Length: "0",
				Type:   "audio/mpeg",
			}
		}

		if ep.ImageURL != "" {
			item.ITunesImage = &ITunesImage{Href: ep.ImageURL}
		}

		if ep.Duration > 0 {
			hours := ep.Duration / 3600
			minutes := (ep.Duration % 3600) / 60
			seconds := ep.Duration % 60
			if hours > 0 {
				item.ITunesDuration = fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
			} else {
				item.ITunesDuration = fmt.Sprintf("%d:%02d", minutes, seconds)
			}
		}

		rss.Channel.Items = append(rss.Channel.Items, item)
	}

	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), output...), nil
}
