package artists

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

func Search(provider *api.Provider, name string) ([]*Artist, error) {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/artists", provider.URL))
	values := u.Query()
	values.Set("name", url.QueryEscape(name))
	u.RawQuery = values.Encode()

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

	artists := []*Artist{}
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func Create(provider *api.Provider, artist *Artist) error {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/artists", provider.URL))

	b, _ := json.Marshal(&artist)
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

	return json.NewDecoder(resp.Body).Decode(&artist)
}

func Get(provider *api.Provider, id int64, _ *GetOptions) (*Artist, error) {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/artists/%d", provider.URL, id))

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

	artist := Artist{}
	if err := json.NewDecoder(resp.Body).Decode(&artist); err != nil {
		return nil, err
	}
	return &artist, nil
}
