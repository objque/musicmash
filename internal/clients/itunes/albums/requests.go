package albums

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/pkg/errors"
)

func GetArtistAlbums(provider *itunes.Provider, artistID uint64) ([]*Album, error) {
	albumsURL := fmt.Sprintf("%s/v1/catalog/us/artists/%v/albums?limit=100", provider.URL, artistID)
	req, _ := http.NewRequest(http.MethodGet, albumsURL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.Token))
	provider.WaitRateLimit()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get albums for %v", artistID)
	}
	defer resp.Body.Close()

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
