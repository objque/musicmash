package db

type Search struct {
	Artists  []*Artist          `json:"artists"`
	Releases []*InternalRelease `json:"releases"`
}

func (mgr *AppDatabaseMgr) Search(query string) (*Search, error) {
	artists, err := mgr.SearchArtists(query)
	if err != nil {
		return nil, err
	}

	rels, err := mgr.GetInternalReleases(&GetInternalReleaseOpts{Title: query})
	if err != nil {
		return nil, err
	}

	result := Search{Artists: artists, Releases: rels}
	return &result, nil
}
