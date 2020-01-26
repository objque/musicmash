package songs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/pkg/errors"
)

func getSongs(provider *itunes.Provider, url string) (*Data, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.Token))
	provider.WaitRateLimit()
	resp, err := provider.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data := Data{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		b, readerErr := ioutil.ReadAll(decoder.Buffered())
		if readerErr != nil {
			return nil, fmt.Errorf("can't read all body because %v", readerErr)
		}
		return nil, errors.Wrapf(err, "looking json, but got %s", string(b))
	}

	return &data, nil
}

func GetArtistSongs(provider *itunes.Provider, artistID uint64) ([]*Song, error) {
	log.Debugf("Getting songs by %d", artistID)
	url := fmt.Sprintf("%s/v1/catalog/us/artists/%v/songs?limit=20", provider.URL, artistID)
	data, err := getSongs(provider, url)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get songs for %v", artistID)
	}

	songs := data.Songs
	for data.Next != "" {
		log.Debugf("Getting next page (%s) of songs for artist %v", data.Next, artistID)
		url := fmt.Sprintf("%s%s&limit=20", provider.URL, data.Next)
		data, err = getSongs(provider, url)
		if err != nil {
			return nil, errors.Wrapf(err, "tried to get songs for %v", artistID)
		}
		songs = append(songs, data.Songs...)
	}

	return songs, nil
}

func isLatest(song *Song) bool {
	now := time.Now().UTC().Truncate(time.Hour * 24)
	yesterday := now.Add(-time.Hour * 48)
	return song.Attributes.ReleaseDate.Value.UTC().After(yesterday)
}

func GetLatestArtistSongs(provider *itunes.Provider, artistID uint64) ([]*Song, error) {
	songs, err := GetArtistSongs(provider, artistID)
	if err != nil {
		return nil, err
	}

	latest := []*Song{}
	for _, release := range songs {
		if !isLatest(release) {
			continue
		}

		latest = append(latest, release)
	}
	return latest, err
}
