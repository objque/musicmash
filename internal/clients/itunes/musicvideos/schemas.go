package musicvideos

import (
	"strconv"
	"strings"
	"time"

	"github.com/musicmash/musicmash/internal/clients/itunes/types"
)

type Data struct {
	MusicVideos []*MusicVideo `json:"data"`
	Next        string        `json:"next"`
}

type MusicVideo struct {
	ID         string               `json:"id"`
	Attributes MusicVideoAttributes `json:"attributes"`
}

type MusicVideoAttributes struct {
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

func (m *MusicVideo) GetPoster(_, _ int) string {
	// todo: width ang height are ignored, because 500x500 is hardcoded, but music-videos with this size wont'open
	width := m.Attributes.Artwork.Width
	height := m.Attributes.Artwork.Height
	url := strings.Replace(m.Attributes.Artwork.URL, "{w}", strconv.Itoa(width), 1)
	return strings.Replace(url, "{h}", strconv.Itoa(height), 1)
}

func (m *MusicVideo) GetID() string {
	return m.ID
}

func (m *MusicVideo) GetName() string {
	return m.Attributes.Name
}

func (m *MusicVideo) GetAlbumName() string {
	return m.Attributes.AlbumName
}

func (m *MusicVideo) GetReleaseDate() time.Time {
	return m.Attributes.ReleaseDate.Value
}
