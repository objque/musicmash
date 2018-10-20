package yandex

import (
	"fmt"
	"strings"

	"github.com/musicmash/musicmash/internal/clients/yandex/types"
)

type Session struct {
	UID string `json:"yandexuid"`
}

type SearchResult struct {
	Artists struct {
		Items []*Artist `json:"items"`
	} `json:"artists"`
}
type Artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ArtistInfo struct {
	Albums []*ArtistAlbum `json:"albums"`
}
type ArtistAlbum struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Released    types.Time `json:"releaseDate"`
	ReleaseYear int        `json:"year"`
	Version     string     `json:"version"`
	Poster      string     `json:"ogImage"`
}

func (a *ArtistAlbum) GetPosterWithSize(width, height int) string {
	return "https://" + strings.Replace(a.Poster, "%%", fmt.Sprintf("%dx%d", width, height), -1)
}
