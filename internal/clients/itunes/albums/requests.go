package albums

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/pkg/errors"
)

func GetArtistAlbums(provider *itunes.Provider, artistID uint64) ([]*Album, error) {
	albumsURL := fmt.Sprintf("%s/v1/catalog/us/artists/%v/albums", provider.URL, artistID)
	req, _ := http.NewRequest(http.MethodGet, albumsURL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.Token))
	provider.WaitRateLimit()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get albums for %v", artistID)
	}

	type answer struct {
		Albums []*Album `json:"data"`
	}
	a := answer{}
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return nil, errors.Wrapf(err, "tried to decode albums for %v", artistID)
	}
	return a.Albums, nil
}

func GetLatestArtistAlbum(provider *itunes.Provider, artistID uint64) (*Album, error) {
	albums, err := GetArtistAlbums(provider, artistID)
	if err != nil {
		return nil, err
	}

	if len(albums) == 0 {
		return nil, ErrAlbumsNotFound
	}

	latest := albums[0]
	for _, album := range albums {
		if album.Attributes.ReleaseDate.Value.After(latest.Attributes.ReleaseDate.Value) {
			latest = album
		}
	}
	return latest, nil
}
