package notifysettings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/musicmash/musicmash/pkg/api"
)

func Create(provider *api.Provider, userName string, settings *Settings) error {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/notifications/settings", provider.URL))

	b, _ := json.Marshal(&settings)
	headers := http.Header{
		"x-user-name": {userName},
	}

	return api.PostWithHeaders(provider, u, headers, bytes.NewBuffer(b), &settings)
}

func List(provider *api.Provider, userName string) ([]*Settings, error) {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/notifications/settings", provider.URL))

	headers := http.Header{
		"x-user-name": {userName},
	}

	settings := []*Settings{}
	if err := api.GetWithHeaders(provider, u, headers, &settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func Update(provider *api.Provider, userName string, settings *Settings) error {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/notifications/settings", provider.URL))

	b, _ := json.Marshal(&settings)
	headers := http.Header{
		"x-user-name": {userName},
	}

	return api.PatchWithHeaders(provider, u, headers, bytes.NewBuffer(b), &settings)
}
