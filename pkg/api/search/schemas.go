package search

import (
	"github.com/musicmash/musicmash/pkg/api/artists"
	"github.com/musicmash/musicmash/pkg/api/releases"
)

type Result struct {
	Artists  []*artists.Artist   `json:"artists"`
	Releases []*releases.Release `json:"releases"`
}
