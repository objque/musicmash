package associations

type Association struct {
	ArtistID  int64  `json:"artist_id"`
	StoreName string `json:"name"`
	StoreID   string `json:"id"`
}

type ListOpts struct {
	ArtistID  int64
	StoreName string
}
