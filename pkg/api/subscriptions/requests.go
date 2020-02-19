package subscriptions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/musicmash/musicmash/pkg/api"
)

func List(provider *api.Provider, userName string) ([]*Subscription, error) {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/subscriptions", provider.URL))

	headers := http.Header{
		"x-user-name": {userName},
	}

	subscriptions := []*Subscription{}
	if err := api.GetWithHeaders(provider, u, headers, &subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func Create(provider *api.Provider, userName string, artists []int64) error {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/subscriptions", provider.URL))

	headers := http.Header{
		"x-user-name": {userName},
	}

	b, _ := json.Marshal(&artists)
	return api.PostWithHeaders(provider, u, headers, bytes.NewBuffer(b), nil)
}

func Delete(provider *api.Provider, userName string, artists []int64) error {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/subscriptions", provider.URL))

	headers := http.Header{
		"x-user-name": {userName},
	}

	b, _ := json.Marshal(&artists)
	return api.DeleteWithHeaders(provider, u, headers, bytes.NewBuffer(b))
}
