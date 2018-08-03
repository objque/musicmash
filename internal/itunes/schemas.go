package itunes

import "time"

type Release struct {
	WrapperType            string    `json:"wrapperType"`
	CollectionType         string    `json:"collectionType"`
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
	ID       uint64
	Date     time.Time
	IsComing bool
}

func NewInfo(id, released string) *LastRelease {
	return &LastRelease{}
}

func (r *LastRelease) IsLatest() bool {
	now := time.Now().UTC().Truncate(time.Hour * 24)
	yesterday := now.Add(-time.Hour * 48)
	return r.Date.UTC().After(yesterday)
}
