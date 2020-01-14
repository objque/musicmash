package artists

import "time"

type Artist struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Poster     string `json:"poster"`
	Popularity int    `json:"popularity"`
	Followers  uint   `json:"followers"`
}

type Association struct {
	ArtistID  int64  `json:"artist_id"`
	StoreName string `json:"name"`
	StoreID   string `json:"id"`
}

type Release struct {
	ID        uint64    `json:"id"`
	ArtistID  int64     `json:"artist_id"`
	Released  time.Time `json:"released"`
	Poster    string    `json:"poster"`
	Title     string    `json:"title"`
	ItunesID  string    `json:"itunes_id"`
	SpotifyID string    `json:"spotify_id"`
	DeezerID  string    `json:"deezer_id"`
}
