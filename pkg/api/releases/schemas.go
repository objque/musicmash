package releases

import "time"

type Release struct {
	ID          uint64    `json:"id"`
	ArtistID    int64     `json:"artist_id"`
	ArtistName  string    `json:"artist_name"`
	Released    time.Time `json:"released"`
	Poster      string    `json:"poster"`
	Title       string    `json:"title"`
	SpotifyID   string    `json:"spotify_id"`
	Type        string    `json:"type"`
	DurationMs  int64     `json:"duration_ms"`
	TracksCount int32     `json:"tracks_count"`
	Explicit    bool      `json:"explicit"`
}
