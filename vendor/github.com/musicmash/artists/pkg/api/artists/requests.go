package artists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/musicmash/artists/pkg/api"
)

func GetFromStore(provider *api.Provider, storeName string) ([]*StoreInfo, error) {
	url := fmt.Sprintf("%s/artists/store/%s", provider.URL, storeName)
	resp, err := provider.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("got %d status code", resp.StatusCode)
	}

	artists := []*StoreInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func Validate(provider *api.Provider, artists []int64) ([]int64, error) {
	body, err := json.Marshal(&artists)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/validate", provider.URL)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := provider.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("got %d status code", resp.StatusCode)
	}

	artists = []int64{}
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func GetFullInfo(provider *api.Provider, ids []int64) ([]*Artist, error) {
	body, err := json.Marshal(&ids)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/artists", provider.URL)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := provider.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("got %d status code", resp.StatusCode)
	}

	artists := []*Artist{}
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}
	return artists, nil
}
