package itunes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/objque/musicmash/internal/config"
	"github.com/pkg/errors"
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

const (
	releasedTag = `<time data-test-we-datetime datetime="`
	IDTag       = `class="featured-album targeted-link"`
)

func decode(buffer []byte) (*LastRelease, error) {
	parts := strings.Split(string(buffer), releasedTag)
	if len(parts) != 2 {
		return nil, errors.New("can't decode: after split by a time-html tag we have not 2 parts")
	}

	// Jul 18, 2018" aria-label="July 18 ...
	released := strings.Split(parts[1], `"`)[0]
	t, err := parseTime(released)
	if err != nil {
		return nil, errors.Wrapf(err, "can't parse time '%s'", released)
	}

	xxx := strings.Split(parts[0], IDTag)[0]
	urls := strings.Split(xxx, `<a href="`)
	return &LastRelease{
		Date: *t,
		URL:  strings.Replace(urls[len(urls)-1], `" `, "", 1),
	}, nil
}

func GetArtistInfo(url string) (*LastRelease, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "can't receive page '%s'", url)
	}

	buffer, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, errors.Wrapf(err, "can't read response '%s'", url)
	}

	info, err := decode(buffer)
	if err != nil {
		return nil, errors.Wrapf(err, "can't decode '%s'", url)
	}
	return info, nil
}
