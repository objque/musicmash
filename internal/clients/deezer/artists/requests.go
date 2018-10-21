package artists

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/musicmash/musicmash/internal/clients/deezer"
	"github.com/pkg/errors"
)

func SearchArtist(provider *deezer.Provider, term string) (*Artist, error) {
	searchURL := fmt.Sprintf("%s/search/artist?limit=1&q=%s", provider.URL, url.QueryEscape(term))
	provider.WaitRateLimit()
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to search artist with name '%v'", term)
	}
	defer resp.Body.Close()

	type answer struct {
		Artists []*Artist `json:"data"`
	}
	a := answer{}
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return nil, errors.Wrapf(err, "tried to decode albums for %v", term)
	}
	if len(a.Artists) == 0 {
		return nil, ErrArtistNotFound
	}
	return a.Artists[0], nil
}
