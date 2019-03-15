package artists

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/pkg/errors"
)

func SearchArtist(provider *itunes.Provider, term string) (*Artist, error) {
	albumsURL := fmt.Sprintf("%s/v1/catalog/us/search?types=artists&limit=1&term=%s", provider.URL, url.QueryEscape(term))
	req, _ := http.NewRequest(http.MethodGet, albumsURL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.Token))
	provider.WaitRateLimit()
	resp, err := provider.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to search artist with name '%v'", term)
	}
	defer resp.Body.Close()

	type answer struct {
		Results struct {
			Data struct {
				Artists []*Artist `json:"data"`
			} `json:"artists"`
		} `json:"results"`
	}
	a := answer{}
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return nil, errors.Wrapf(err, "tried to decode search result for artist with name '%v'", term)
	}
	if len(a.Results.Data.Artists) == 0 {
		return nil, ErrArtistNotFound
	}
	return a.Results.Data.Artists[0], nil
}
