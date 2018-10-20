package albums

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/musicmash/musicmash/internal/clients/deezer"
	"github.com/pkg/errors"
)

func GetArtistAlbums(provider *deezer.Provider, artistID int) ([]*Album, error) {
	albumsURL := fmt.Sprintf("%s/artist/%d/albums?limit=50", provider.URL, artistID)
	resp, err := http.Get(albumsURL)
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

func GetLatestArtistAlbum(provider *deezer.Provider, artistID int) (*Album, error) {
	albums, err := GetArtistAlbums(provider, artistID)
	if err != nil {
		return nil, err
	}

	if len(albums) == 0 {
		return nil, ErrAlbumsNotFound
	}

	latest := albums[0]
	for _, album := range albums {
		if album.Released.Value.After(latest.Released.Value) {
			latest = album
		}
	}
	return latest, nil
}
