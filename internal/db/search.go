package db

type Search struct {
	Artists []*Artist `json:"artists"`
}

type SearchMgr interface {
	Search(query string) (*Search, error)
}

func (mgr *AppDatabaseMgr) Search(query string) (*Search, error) {
	artists, err := mgr.SearchArtists(query)
	if err != nil {
		return nil, err
	}

	result := Search{Artists: artists}
	return &result, nil
}
