package artists

type Artist struct {
	ID         string           `json:"id"`
	Attributes ArtistAttributes `json:"attributes"`
}

type ArtistAttributes struct {
	Name string `json:"name"`
}
