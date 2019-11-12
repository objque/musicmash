package albums

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/musicmash/musicmash/internal/clients/deezer"
	"github.com/pkg/errors"
)

func GetArtistAlbums(provider *deezer.Provider, artistID int) ([]*Album, error) {
	albumsURL := fmt.Sprintf("%s/artist/%d/albums?limit=50", provider.URL, artistID)
	provider.WaitRateLimit()
	resp, err := http.Get(albumsURL)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get albums for %v", artistID)
	}
	defer func() { _ = resp.Body.Close() }()

	type answer struct {
		Albums []*Album `json:"data"`
	}
	a := answer{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&a); err != nil {
		errBody, err := ioutil.ReadAll(dec.Buffered())
		if err != nil {
			return nil, errors.Wrapf(err, "tried to decode albums for %v", artistID)
		}
		return nil, errors.Wrapf(err, "tried to decode albums for %v from: %s", artistID, string(errBody))
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

func GetByID(provider *deezer.Provider, albumID int) (*Album, error) {
	albumURL := fmt.Sprintf("%s/album/%d", provider.URL, albumID)
	provider.WaitRateLimit()
	resp, err := http.Get(albumURL)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get album with id %v", albumID)
	}
	defer func() { _ = resp.Body.Close() }()

	album := Album{}
	if err := json.NewDecoder(resp.Body).Decode(&album); err != nil {
		return nil, errors.Wrapf(err, "tried to decode album with id %v", albumID)
	}
	return &album, nil
}
