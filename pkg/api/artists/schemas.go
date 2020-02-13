package artists

type Artist struct {
	ID         int64  `json:"id,omitempty"`
	Name       string `json:"name"`
	Poster     string `json:"poster"`
	Popularity int    `json:"popularity"`
	Followers  uint   `json:"followers"`
}
