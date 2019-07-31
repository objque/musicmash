package artists

type StoreInfo struct {
	ArtistID  int64  `json:"artist_id"`
	StoreName string `json:"name"`
	StoreID   string `json:"id"`
}

type Artist struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Poster     string `json:"poster"`
	Popularity int    `json:"popularity"`
	Followers  uint   `json:"followers"`
}
