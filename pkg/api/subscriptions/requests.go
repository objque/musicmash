package subscriptions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/musicmash/musicmash/pkg/api"
	log "github.com/sirupsen/logrus"
	"moul.io/http2curl"
)

func List(provider *api.Provider, userName string) ([]*Subscription, error) {
	url := fmt.Sprintf("%s/subscriptions", provider.URL)
	request, err := http.NewRequest(http.MethodGet, url, nil)
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

	artists := []*Subscription{}
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func Create(provider *api.Provider, userName string, artists []int64) error {
	url := fmt.Sprintf("%s/subscriptions", provider.URL)
	b, _ := json.Marshal(&artists)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
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

	if resp.StatusCode == http.StatusCreated {
		return nil
	}
	return api.ExtractError(resp.Body)
}

func Delete(provider *api.Provider, userName string, artists []int64) error {
	url := fmt.Sprintf("%s/subscriptions", provider.URL)
	b, _ := json.Marshal(&artists)
	request, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(b))
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

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	return api.ExtractError(resp.Body)
}
