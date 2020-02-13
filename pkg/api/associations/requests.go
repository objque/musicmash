package associations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/musicmash/musicmash/pkg/api"
	log "github.com/sirupsen/logrus"
	"moul.io/http2curl"
)

func Create(provider *api.Provider, info *Association) error {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/associations", provider.URL))
	b, _ := json.Marshal(&info)
	request, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	command, _ := http2curl.GetCurlCommand(request)
	log.Debug(command)

	resp, err := provider.Client.Do(request)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return api.ExtractError(resp.Body)
	}

	return json.NewDecoder(resp.Body).Decode(&info)
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

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	command, _ := http2curl.GetCurlCommand(request)
	log.Debug(command)

	resp, err := provider.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, api.ExtractError(resp.Body)
	}

	associations := []*Association{}
	if err := json.NewDecoder(resp.Body).Decode(&associations); err != nil {
		return nil, err
	}
	return associations, nil
}
