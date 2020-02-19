package associations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/musicmash/musicmash/pkg/api"
)

func Create(provider *api.Provider, info *Association) error {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/associations", provider.URL))

	b, _ := json.Marshal(info)
	return api.Post(provider, u, bytes.NewBuffer(b), &info)
}

func List(provider *api.Provider, opts *ListOpts) ([]*Association, error) {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/associations", provider.URL))
	if opts != nil {
		values := u.Query()
		if opts.ArtistID != 0 {
			values.Set("artist_id", fmt.Sprint(opts.ArtistID))
		}
		if opts.StoreName != "" {
			values.Set("store_name", opts.StoreName)
		}
		u.RawQuery = values.Encode()
	}

	associations := []*Association{}
	if err := api.Get(provider, u, &associations); err != nil {
		return nil, err
	}

	return associations, nil
}
