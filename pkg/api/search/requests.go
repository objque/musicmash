package search

import (
	"fmt"
	"net/url"

	"github.com/musicmash/musicmash/pkg/api"
)

func Do(provider *api.Provider, query string) (*Result, error) {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/search", provider.URL))

	values := u.Query()
	values.Set("query", url.QueryEscape(query))
	u.RawQuery = values.Encode()

	result := Result{}
	if err := api.Get(provider, u, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
