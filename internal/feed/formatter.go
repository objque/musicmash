package feed

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/feeds"
	"github.com/musicmash/musicmash/internal/db"
)

type Formatter struct {
	Title       string
	Link        string
	Description string
}

func makeRssDescription(release *db.Release) string {
	return fmt.Sprintf("[%s] %s –– %s", release.StoreName, release.ArtistName, release.Title)
}

func (f *Formatter) ToRawRss(releases []*db.Release) *feeds.Feed {
	result := &feeds.Feed{
		Title:       f.Title,
		Link:        &feeds.Link{Href: f.Link},
		Description: f.Description,
	}

	result.Items = make([]*feeds.Item, len(releases))
	for i, release := range releases {
		title := makeRssDescription(release)
		result.Items[i] = &feeds.Item{
			Title:       title,
			Description: title,
			Created:     release.Released,
		}
		// TODO (m.kalinin): extract links into config
		switch release.StoreName {
		case "deezer":
			result.Items[i].Link = &feeds.Link{Href: fmt.Sprintf("https://www.deezer.com/en/album/%s", release.StoreID)}
		case "itunes":
			result.Items[i].Link = &feeds.Link{Href: fmt.Sprintf("https://itunes.apple.com/us/album/%s", release.StoreID)}
		}
	}
	return result
}

func (f *Formatter) ToRss(releases []*db.Release) ([]byte, error) {
	result := f.ToRawRss(releases)
	b, err := result.ToRss()
	if err != nil {
		return nil, err
	}
	return []byte(b), nil
}

func (f *Formatter) ToJson(releases []*db.Release) ([]byte, error) {
	return json.Marshal(&releases)
}
