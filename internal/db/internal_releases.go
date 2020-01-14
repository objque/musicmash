package db

import (
	"time"
)

type InternalRelease struct {
	ID        uint64    `json:"id"`
	ArtistID  int64     `json:"artist_id"`
	Released  time.Time `json:"released"`
	Poster    string    `json:"poster"`
	Title     string    `json:"title"`
	ItunesID  string    `json:"itunes_id"`
	SpotifyID string    `json:"spotify_id"`
	DeezerID  string    `json:"deezer_id"`
}

type InternalReleaseMgr interface {
	GetArtistInternalReleases(id int64) ([]*InternalRelease, error)
	GetArtistsInternalReleases(ids []int64) ([]*InternalRelease, error)
}

func (mgr *AppDatabaseMgr) GetArtistInternalReleases(id int64) ([]*InternalRelease, error) {
	return mgr.GetArtistsInternalReleases([]int64{id})
}

func (mgr *AppDatabaseMgr) GetArtistsInternalReleases(ids []int64) ([]*InternalRelease, error) {
	const query = `
SELECT releases.id,
       releases.artist_id,
	   releases.released,
	   releases.poster,
	   releases.title,
	   itunes.store_id  AS itunes_id,
	   spotify.store_id AS spotify_id,
	   deezer.store_id  AS deezer_id
FROM releases AS releases
LEFT JOIN releases AS itunes ON (
   releases.artist_id = itunes.artist_id AND
   releases.title     = itunes.title     AND
   itunes.store_name  = 'itunes'
)
LEFT JOIN releases AS spotify ON (
   releases.artist_id = spotify.artist_id AND
   releases.title     = spotify.title     AND
   spotify.store_name = 'spotify'
)
LEFT JOIN releases AS deezer ON (
   releases.artist_id = deezer.artist_id AND
   releases.title     = deezer.title     AND
   deezer.store_name  = 'deezer'
)
WHERE releases.artist_id in (?)
GROUP BY releases.title
ORDER BY releases.released DESC
`
	releases := []*InternalRelease{}
	err := mgr.db.Raw(query, ids).Scan(&releases).Error
	if err != nil {
		return nil, err
	}
	return releases, nil
}
