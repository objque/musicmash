package albums

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

func getAlbums(provider *itunes.Provider, url string) (*Data, error) {
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

func GetArtistAlbums(provider *itunes.Provider, artistID uint64) ([]*Album, error) {
	albumsURL := fmt.Sprintf("%s/v1/catalog/us/artists/%v/albums?limit=100", provider.URL, artistID)
	data, err := getAlbums(provider, albumsURL)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get albums for %v", artistID)
	}

	albums := data.Albums
	for data.Next != "" {
		log.Debugf("Getting next page (%s) for artist %v", data.Next, artistID)
		url := fmt.Sprintf("%s%s&limit=100", provider.URL, data.Next)
		data, err = getAlbums(provider, url)
		if err != nil {
			return nil, errors.Wrapf(err, "tried to get albums for %v", artistID)
		}
		albums = append(albums, data.Albums...)
	}

	return albums, nil
}

func isLatest(album *Album) bool {
	now := time.Now().UTC().Truncate(time.Hour * 24)
	yesterday := now.Add(-time.Hour * 48)
	return album.Attributes.ReleaseDate.Value.UTC().After(yesterday)
}

func GetLatestArtistAlbums(provider *itunes.Provider, artistID uint64) ([]*Album, error) {
	albums, err := GetArtistAlbums(provider, artistID)
	if err != nil {
		return nil, err
	}

	latest := []*Album{}
	for _, release := range albums {
		if !isLatest(release) {
			continue
		}

		latest = append(latest, release)
	}
	return latest, err
}
