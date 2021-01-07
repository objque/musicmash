package releases

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/musicmash/musicmash/pkg/api"
)

const (
	SortTypeASC  = "asc"
	SortTypeDESC = "desc"
)

type GetOptions struct {
	Before      *uint64
	Limit       *uint64
	Offset      *uint64
	ArtistID    *int64
	UserName    string
	ReleaseType string
	SortType    string
	Explicit    *bool
	Since       *time.Time
	Till        *time.Time
}

func buildValues(opts *GetOptions) *url.Values {
	values := url.Values{}

	if opts.Before != nil {
		values.Set("before", fmt.Sprintf("%v", *opts.Before))
	}

	if opts.Limit != nil {
		values.Set("limit", fmt.Sprintf("%v", *opts.Limit))
	}

	if opts.Offset != nil {
		values.Set("offset", fmt.Sprintf("%v", *opts.Offset))
	}

	if opts.ArtistID != nil {
		values.Set("artist_id", fmt.Sprintf("%v", *opts.ArtistID))
	}

	if opts.ReleaseType != "" {
		values.Set("type", opts.ReleaseType)
	}

	if opts.Explicit != nil {
		values.Set("explicit", fmt.Sprintf("%v", *opts.Explicit))
	}

	if opts.SortType != "" {
		values.Set("sort_type", opts.SortType)
	}

	if opts.Since != nil {
		values.Set("since", opts.Since.Format("2006-01-02"))
	}

	if opts.Till != nil {
		values.Set("till", opts.Till.Format("2006-01-02"))
	}

	return &values
}

func List(provider *api.Provider, opts *GetOptions) ([]*Release, error) {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/releases", provider.URL))
	headers := http.Header{}
	if opts != nil {
		u.RawQuery = buildValues(opts).Encode()

		if opts.UserName != "" {
			headers.Set("x-user-name", opts.UserName)
		}
	}

	releases := []*Release{}
	if err := api.GetWithHeaders(provider, u, headers, &releases); err != nil {
		return nil, err
	}

	return releases, nil
}
