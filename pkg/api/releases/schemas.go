package releases

import "time"

type Release struct {
	ID         uint64    `json:"id"`
	ArtistID   int64     `json:"artist_id"`
	ArtistName string    `json:"artist_name"`
	Released   time.Time `json:"released"`
	Poster     string    `json:"poster"`
	Title      string    `json:"title"`
	ItunesID   *string   `json:"itunes_id"`
	SpotifyID  *string   `json:"spotify_id"`
	DeezerID   *string   `json:"deezer_id"`
	Type       string    `json:"type"`
	Explicit   bool      `json:"explicit"`
}
