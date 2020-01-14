package stores

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/musicmash/musicmash/pkg/api"
	log "github.com/sirupsen/logrus"
	"moul.io/http2curl"
)

func List(provider *api.Provider) ([]*Store, error) {
	url := fmt.Sprintf("%s/stores", provider.URL)
	request, err := http.NewRequest(http.MethodGet, url, nil)
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

	stores := []*Store{}
	if err := json.NewDecoder(resp.Body).Decode(&stores); err != nil {
		return nil, err
	}
	return stores, nil
}
