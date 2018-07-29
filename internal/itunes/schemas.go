package itunes

import "time"

type Release struct {
	WrapperType            string    `json:"wrapperType"`
	Kind                   string    `json:"kind"`
	ArtistID               int       `json:"artistId"`
	CollectionID           int       `json:"collectionId"`
	TrackID                uint64    `json:"trackId"`
	ArtistName             string    `json:"artistName"`
	CollectionName         string    `json:"collectionName"`
	TrackName              string    `json:"trackName"`
	CollectionCensoredName string    `json:"collectionCensoredName"`
	TrackCensoredName      string    `json:"trackCensoredName"`
	ArtistViewURL          string    `json:"artistViewUrl"`
	CollectionViewURL      string    `json:"collectionViewUrl"`
	TrackViewURL           string    `json:"trackViewUrl"`
	PreviewURL             string    `json:"previewUrl"`
	ArtworkURL30           string    `json:"artworkUrl30"`
	ArtworkURL60           string    `json:"artworkUrl60"`
	ArtworkURL100          string    `json:"artworkUrl100"`
	CollectionPrice        float64   `json:"collectionPrice"`
	TrackPrice             float64   `json:"trackPrice"`
	ReleaseDate            time.Time `json:"releaseDate"`
	CollectionExplicitness string    `json:"collectionExplicitness"`
	TrackExplicitness      string    `json:"trackExplicitness"`
	DiscCount              int       `json:"discCount"`
	DiscNumber             int       `json:"discNumber"`
	TrackCount             int       `json:"trackCount"`
	TrackNumber            int       `json:"trackNumber"`
	TrackTimeMillis        int       `json:"trackTimeMillis"`
	Country                string    `json:"country"`
	Currency               string    `json:"currency"`
	PrimaryGenreName       string    `json:"primaryGenreName"`
	IsStreamable           bool      `json:"isStreamable"`
	ContentAdvisoryRating  string    `json:"contentAdvisoryRating,omitempty"`
	CollectionArtistID     int       `json:"collectionArtistId,omitempty"`
	CollectionArtistName   string    `json:"collectionArtistName,omitempty"`
}

type SearchResult struct {
	Releases      []*Release `json:"results"`
	ReleasesCount int        `json:"resultCount"`
}

func (r *Release) IsLatest() bool {
	now := time.Now().UTC().Truncate(time.Hour * 24)
	yesterday := now.Add(-time.Hour * 48)
	return r.ReleaseDate.UTC().After(yesterday)
}

type LastRelease struct {
	URL  string
	Date time.Time
}

func NewInfo(id, released string) *LastRelease {
	return &LastRelease{}
}

func (r *LastRelease) IsLatest() bool {
	now := time.Now().UTC().Truncate(time.Hour * 24)
	yesterday := now.Add(-time.Hour * 48)
	return r.Date.UTC().After(yesterday)
}
