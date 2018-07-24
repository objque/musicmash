package itunes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/objque/musicmash/internal/config"
)

const (
	album = "album"
	track = "musicTrack"
)

func findLastRelease(releases []*Release) *Release {
	last := releases[0]
	for _, release := range releases {
		if release.ReleaseDate.After(last.ReleaseDate) {
			last = release
			continue
		}
	}
	return last
}

func get(artist, entity string) (*Release, error) {
	resp, err := http.Get(fmt.Sprintf("%s/search?term=%s&media=music&entity=%s", config.Config.StoreURL, url.QueryEscape(artist), entity))
	if err != nil {
		return nil, err
	}

	result := SearchResult{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	resp.Body.Close()

	if result.ReleasesCount == 0 {
		return nil, fmt.Errorf("artist '%s' not found", artist)
	}

	lastRelease := findLastRelease(result.Releases)
	if !strings.Contains(strings.ToLower(lastRelease.ArtistName), strings.ToLower(artist)) {
		return nil, fmt.Errorf("founded artist '%s' not equals to '%s'", lastRelease.ArtistName, artist)
	}
	return lastRelease, nil
}

func GetLatestTrackRelease(artist string) (*Release, error) {
	return get(artist, track)
}

func GetLatestAlbumRelease(artist string) (*Release, error) {

	return get(artist, album)
}
