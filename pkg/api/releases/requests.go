package releases

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/musicmash/musicmash/pkg/api"
)

func Get(provider *api.Provider, since time.Time) ([]*Release, error) {
	url := fmt.Sprintf("%s/releases?since=%s", provider.URL, since.Format("2006-01-02T15:04:05"))
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

	release := []*Release{}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return release, nil
}
