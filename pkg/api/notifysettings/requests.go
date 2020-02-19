package notifysettings

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

func Create(provider *api.Provider, userName string, settings *Settings) error {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/notifications/settings", provider.URL))

	b, _ := json.Marshal(&settings)
	request, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	request.Header.Add("x-user-name", userName)

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

	return json.NewDecoder(resp.Body).Decode(&settings)
}

func List(provider *api.Provider, userName string) ([]*Settings, error) {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/notifications/settings", provider.URL))

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("x-user-name", userName)

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

	settings := []*Settings{}
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return nil, err
	}
	return settings, nil
}

func Update(provider *api.Provider, userName string, settings *Settings) error {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/notifications/settings", provider.URL))

	b, _ := json.Marshal(&settings)
	request, err := http.NewRequest(http.MethodPatch, u.String(), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	request.Header.Add("x-user-name", userName)

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

	return json.NewDecoder(resp.Body).Decode(&settings)
}
