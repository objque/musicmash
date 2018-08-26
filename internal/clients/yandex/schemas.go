package yandex

type Session struct {
	UID string `json:"yandexuid"`
}

type SearchResult struct {
	Artists struct {
		Items []*Artist `json:"items"`
	} `json:"artists"`
}
type Artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ArtistInfo struct {
	Albums []*ArtistAlbum `json:"albums"`
}
type ArtistAlbum struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Released    string `json:"releaseDate"`
	ReleaseYear int    `json:"year"`
	Version     string `json:"version"`
}
