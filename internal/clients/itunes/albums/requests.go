package albums

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/pkg/errors"
)

func getAlbums(provider *itunes.Provider, url string) (*Data, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.Token))
	provider.WaitRateLimit()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := Data{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func GetArtistAlbums(provider *itunes.Provider, artistID uint64) ([]*Album, error) {
	albumsURL := fmt.Sprintf("%s/v1/catalog/us/artists/%v/albums?limit=100", provider.URL, artistID)
	data, err := getAlbums(provider, albumsURL)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get albums for %v", artistID)
	}

	return data.Albums, nil
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
