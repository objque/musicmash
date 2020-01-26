package songs

import (
	"strconv"
	"strings"
	"time"

	"github.com/musicmash/musicmash/internal/clients/itunes/types"
)

type Data struct {
	Songs []*Song `json:"data"`
	Next  string  `json:"next"`
}

type Song struct {
	ID         string         `json:"id"`
	Attributes SongAttributes `json:"attributes"`
}

type SongAttributes struct {
	Artwork     Artwork    `json:"artwork"`
	ReleaseDate types.Time `json:"releaseDate"`
	Name        string     `json:"name"`
	AlbumName   string     `json:"albumName"`
	URL         string     `json:"url"`
}

type Artwork struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

func (a *Song) GetPoster(width, height int) string {
	url := strings.Replace(a.Attributes.Artwork.URL, "{w}", strconv.Itoa(width), 1)
	return strings.Replace(url, "{h}", strconv.Itoa(height), 1)
}

func (a *Song) GetID() string {
	// https://music.apple.com/us/album/burial-feat-pusha-t-moody-good-trollphace/1015156632?i=1015156895
	// -> 1015156632?i=1015156895 (album_id?song_id)
	args := strings.Split(a.Attributes.URL, "/")
	return args[len(args)-1]
}

func (a *Song) GetName() string {
	return a.Attributes.Name
}

func (a *Song) GetAlbumName() string {
	return a.Attributes.AlbumName
}

func (a *Song) GetReleaseDate() time.Time {
	return a.Attributes.ReleaseDate.Value
}
