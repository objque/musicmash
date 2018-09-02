package albums

import (
	"github.com/objque/musicmash/internal/clients/itunes/v2/types"
)

type Album struct {
	ID         string          `json:"id"`
	Attributes AlbumAttributes `json:"attributes"`
}

type AlbumAttributes struct {
	Name        string     `json:"name"`
	ReleaseDate types.Time `json:"releaseDate"`
	ArtistName  string     `json:"artistName"`
	IsSingle    bool       `json:"isSingle"`
	IsComplete  bool       `json:"isComplete"`
}

type Song struct {
	ID         string         `json:"id"`
	Attributes SongAttributes `json:"attributes"`
}

type SongAttributes struct {
	Name        string        `json:"name"`
	ReleaseDate string        `json:"releaseDate"`
	Date        types.Time    `json:"date"`
	ArtistName  string        `json:"artistName"`
	Previews    []interface{} `json:"previews"`
}

func (s *Song) IsAvailable() bool {
	return len(s.Attributes.Previews) > 0
}
