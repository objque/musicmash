package stores

import (
	"fmt"
	"net/url"

	"github.com/musicmash/musicmash/pkg/api"
)

func List(provider *api.Provider) ([]*Store, error) {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/stores", provider.URL))

	stores := []*Store{}
	if err := api.Get(provider, u, &stores); err != nil {
		return nil, err
	}

	return stores, nil
}
