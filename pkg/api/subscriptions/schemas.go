package subscriptions

type Subscription struct {
	ID           uint64 `json:"id"`
	ArtistID     int64  `json:"artist_id"`
	ArtistName   string `json:"artist_name"`
	ArtistPoster string `json:"artist_poster"`
}
