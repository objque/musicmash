package releases

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/musicmash/musicmash/pkg/api"
)

type GetOptions struct {
	Since, Till *time.Time
}

func For(provider *api.Provider, userName string, opts *GetOptions) ([]*Release, error) {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/releases", provider.URL))
	if opts != nil {
		values := u.Query()
		if opts.Since != nil {
			values.Set("since", opts.Since.Format("2006-01-02"))
		}
		if opts.Till != nil {
			values.Set("till", opts.Till.Format("2006-01-02"))
		}
		u.RawQuery = values.Encode()
	}

	headers := http.Header{
		"x-user-name": {userName},
	}

	releases := []*Release{}
	if err := api.GetWithHeaders(provider, u, headers, &releases); err != nil {
		return nil, err
	}

	return releases, nil
}

func By(provider *api.Provider, id int64) ([]*Release, error) {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/releases", provider.URL))
	values := u.Query()
	values.Set("artist_id", fmt.Sprint(id))
	u.RawQuery = values.Encode()

	releases := []*Release{}
	if err := api.Get(provider, u, &releases); err != nil {
		return nil, err
	}

	return releases, nil
}
