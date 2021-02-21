package artists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/musicmash/musicmash/pkg/api"
)

func Create(provider *api.Provider, artist *Artist) error {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/artists", provider.URL))

	b, _ := json.Marshal(artist)
	return api.Post(provider, u, bytes.NewBuffer(b), &artist)
}

func Get(provider *api.Provider, id int64, _ *GetOptions) (*Artist, error) {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/artists/%d", provider.URL, id))

	artist := &Artist{}
	if err := api.Get(provider, u, artist); err != nil {
		return nil, err
	}

	return artist, nil
}
