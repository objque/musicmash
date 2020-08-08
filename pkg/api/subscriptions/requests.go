package subscriptions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/musicmash/musicmash/pkg/api"
)

type GetOptions struct {
	Limit    *uint64
	Offset   *uint64
	SortType string
}

func buildValues(opts *GetOptions) *url.Values {
	values := url.Values{}

	if opts.SortType != "" {
		values.Set("sort_type", opts.SortType)
	}

	if opts.Limit != nil {
		values.Set("limit", fmt.Sprintf("%v", *opts.Limit))
	}

	if opts.Offset != nil {
		values.Set("offset", fmt.Sprintf("%v", *opts.Offset))
	}

	return &values
}

func List(provider *api.Provider, userName string, opts *GetOptions) ([]*Subscription, error) {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/subscriptions", provider.URL))
	headers := http.Header{
		"x-user-name": {userName},
	}
	if opts != nil {
		u.RawQuery = buildValues(opts).Encode()
	}

	subscriptions := []*Subscription{}
	if err := api.GetWithHeaders(provider, u, headers, &subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func Create(provider *api.Provider, userName string, subscriptions []*Subscription) error {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/subscriptions", provider.URL))

	headers := http.Header{
		"x-user-name": {userName},
	}

	b, _ := json.Marshal(&subscriptions)
	return api.PostWithHeaders(provider, u, headers, bytes.NewBuffer(b), nil)
}

func Delete(provider *api.Provider, userName string, subscriptions []*Subscription) error {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/subscriptions", provider.URL))

	headers := http.Header{
		"x-user-name": {userName},
	}

	b, _ := json.Marshal(&subscriptions)
	return api.DeleteWithHeaders(provider, u, headers, bytes.NewBuffer(b))
}
