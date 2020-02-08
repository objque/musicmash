package albums

import (
	"strconv"
	"strings"
	"time"

	"github.com/musicmash/musicmash/internal/clients/itunes/types"
)

const (
	AlbumReleaseType  = "Album"
	SingleReleaseType = "Single"
	EPReleaseType     = "EP"
	LPReleaseType     = "LP"

	SingleReleaseTypePattern = "- single"
	EPReleaseTypePattern     = " ep"
	LPReleaseTypePattern     = " lp"
)

type Data struct {
	Albums []*Album `json:"data"`
	Next   string   `json:"next,omitempty"`
}

type Album struct {
	ID         string          `json:"id"`
	Attributes AlbumAttributes `json:"attributes"`
}

type AlbumAttributes struct {
	Name          string        `json:"name"`
	ReleaseDate   types.Time    `json:"releaseDate"`
	ArtistName    string        `json:"artistName"`
	IsSingle      bool          `json:"isSingle"`
	IsComplete    bool          `json:"isComplete"`
	Artwork       *AlbumArtwork `json:"artwork"`
	TrackCount    int           `json:"trackCount"`
	ContentRating string        `json:"contentRating"`
}

type AlbumArtwork struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

func (a *AlbumAttributes) GetCollectionType() string {
	// NOTE (m.kalinin): this article is really useful
	// https://support.tunecore.com/hc/en-us/articles/115006689928
	title := strings.ToLower(a.Name)
	if strings.Contains(title, SingleReleaseTypePattern) {
		return SingleReleaseType
	}

	// NOTE (m.kalinin): some releases have < 3 chars
	// for example: 1432623274
	if len(title) > 3 {
		// NOTE (m.kalinin): some EP's ends with Ep instead of '- EP'
		// for example album_id: 1380811617
		if title[len(title)-3:] == EPReleaseTypePattern {
			return EPReleaseType
		}

		// NOTE (m.kalinin): sometimes we have LP
		// for example album_id: 1363601736
		if title[len(title)-3:] == LPReleaseTypePattern {
			return LPReleaseType
		}
	}

	switch {
	case a.TrackCount <= 3:
		return SingleReleaseType
	case a.TrackCount <= 6:
		return EPReleaseType
	default:
		return AlbumReleaseType
	}
}

func (a *Album) GetPoster(width, height int) string {
	url := strings.Replace(a.Attributes.Artwork.URL, "{w}", strconv.Itoa(width), 1)
	return strings.Replace(url, "{h}", strconv.Itoa(height), 1)
}

func (a *Album) GetID() string {
	return a.ID
}

func (a *Album) GetName() string {
	return a.Attributes.Name
}

func (a *Album) GetReleaseDate() time.Time {
	return a.Attributes.ReleaseDate.Value
}

func (a *Album) IsExplicit() bool {
	return a.Attributes.ContentRating == "explicit"
}
