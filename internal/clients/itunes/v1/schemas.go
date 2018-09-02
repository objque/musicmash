package v1

import (
	"strings"
	"time"
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

type Release struct {
	WrapperType            string    `json:"wrapperType"`
	ArtistID               int       `json:"artistId"`
	CollectionID           int       `json:"collectionId"`
	AmgArtistID            int       `json:"amgArtistId"`
	ArtistName             string    `json:"artistName"`
	CollectionName         string    `json:"collectionName"`
	CollectionCensoredName string    `json:"collectionCensoredName"`
	ArtistViewURL          string    `json:"artistViewUrl"`
	CollectionViewURL      string    `json:"collectionViewUrl"`
	ArtworkURL60           string    `json:"artworkUrl60"`
	ArtworkURL100          string    `json:"artworkUrl100"`
	CollectionPrice        float64   `json:"collectionPrice"`
	CollectionExplicitness string    `json:"collectionExplicitness"`
	TrackCount             int       `json:"trackCount"`
	Copyright              string    `json:"copyright"`
	Country                string    `json:"country"`
	Currency               string    `json:"currency"`
	Released               time.Time `json:"releaseDate"`
	PrimaryGenreName       string    `json:"primaryGenreName"`
}

type SearchReleaseResponse struct {
	Count   int        `json:"resultCount"`
	Results []*Release `json:"results"`
}

type LastRelease struct {
	ID         uint64
	Date       time.Time
	IsComing   bool
	ArtistName string
}

type Artist struct {
	Name    string
	StoreID uint64
}

func NewInfo(id, released string) *LastRelease {
	return &LastRelease{}
}

func (r *LastRelease) IsLatest() bool {
	now := time.Now().UTC().Truncate(time.Hour * 24)
	yesterday := now.Add(-time.Hour * 48)
	return r.Date.UTC().After(yesterday)
}

func (r *Release) GetCollectionType() string {
	// NOTE (m.kalinin): this article is really useful
	// https://support.tunecore.com/hc/en-us/articles/115006689928
	title := strings.ToLower(r.CollectionName)
	if strings.Contains(title, SingleReleaseTypePattern) {
		return SingleReleaseType
	}

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

	switch {
	case r.TrackCount <= 3:
		return SingleReleaseType
	case r.TrackCount <= 6:
		return EPReleaseType
	default:
		return AlbumReleaseType
	}
}
