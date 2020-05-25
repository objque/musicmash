package search

import "github.com/musicmash/musicmash/pkg/api/artists"

type Result struct {
	Artists []*artists.Artist `json:"artists"`
}
