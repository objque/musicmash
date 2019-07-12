package releases

import "time"

type Release struct {
	ArtistID  int64     `json:"artist_id"`
	Title     string    `json:"title"`
	Poster    string    `json:"poster"`
	Released  time.Time `json:"released"`
	StoreName string    `json:"store_name"`
	StoreID   string    `json:"store_id"`
}
