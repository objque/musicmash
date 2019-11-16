package notifysettings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/musicmash/musicmash/pkg/api"
)

func Create(provider *api.Provider, userName string, settings *Settings) error {
	url := fmt.Sprintf("%s/notifications/settings", provider.URL)
	b, _ := json.Marshal(&settings)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	request.Header.Add("x-user-name", userName)

	resp, err := provider.Client.Do(request)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return api.ExtractError(resp.Body)
	}

	return json.NewDecoder(resp.Body).Decode(&settings)
}

func List(provider *api.Provider, userName string) ([]*Settings, error) {
	url := fmt.Sprintf("%s/notifications/settings", provider.URL)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("x-user-name", userName)

	resp, err := provider.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, api.ExtractError(resp.Body)
	}

	settings := []*Settings{}
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return nil, err
	}
	return settings, nil
}

func Update(provider *api.Provider, userName string, settings *Settings) error {
	url := fmt.Sprintf("%s/notifications/settings", provider.URL)
	b, _ := json.Marshal(&settings)
	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	request.Header.Add("x-user-name", userName)

	resp, err := provider.Client.Do(request)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return api.ExtractError(resp.Body)
	}

	return json.NewDecoder(resp.Body).Decode(&settings)
}
