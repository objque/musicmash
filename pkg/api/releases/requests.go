package releases

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/musicmash/musicmash/pkg/api"
	log "github.com/sirupsen/logrus"
	"moul.io/http2curl"
)

type GetOptions struct {
	Since, Till *time.Time
}

func For(provider *api.Provider, userName string, opts *GetOptions) ([]*Release, error) {
	url := fmt.Sprintf("%s/releases", provider.URL)
	if opts != nil {
		args := []string{}
		if opts.Since != nil {
			args = append(args, fmt.Sprintf("since=%s", opts.Since.Format("2006-01-02")))
		}
		if opts.Till != nil {
			args = append(args, fmt.Sprintf("till=%s", opts.Till.Format("2006-01-02")))
		}
		url = fmt.Sprintf("%v?%v", url, strings.Join(args, "&"))
	}
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

	stores := []*Release{}
	if err := json.NewDecoder(resp.Body).Decode(&stores); err != nil {
		return nil, err
	}
	return stores, nil
}

func By(provider *api.Provider, id int64) ([]*Release, error) {
	url := fmt.Sprintf("%s/artists/%d/releases", provider.URL, id)
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

	release := []*Release{}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return release, nil
}
