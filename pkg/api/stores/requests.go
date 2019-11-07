package stores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/musicmash/musicmash/pkg/api"
)

func Create(provider *api.Provider, name string) (*Store, error) {
	url := fmt.Sprintf("%s/stores", provider.URL)
	b, _ := json.Marshal(&Store{Name: name})
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	resp, err := provider.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, api.ExtractError(resp.Body)
	}

	store := Store{}
	if err := json.NewDecoder(resp.Body).Decode(&store); err != nil {
		return nil, err
	}
	return &store, nil
}

func List(provider *api.Provider) ([]*Store, error) {
	url := fmt.Sprintf("%s/stores", provider.URL)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := provider.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, api.ExtractError(resp.Body)
	}

	stores := []*Store{}
	if err := json.NewDecoder(resp.Body).Decode(&stores); err != nil {
		return nil, err
	}
	return stores, nil
}
