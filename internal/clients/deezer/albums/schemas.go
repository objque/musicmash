package albums

import (
	"github.com/musicmash/musicmash/internal/clients/deezer/types"
)

type Album struct {
	ID       int        `json:"id"`
	Title    string     `json:"title"`
	Poster   string     `json:"cover_big"`
	Released types.Time `json:"release_date"`
}
