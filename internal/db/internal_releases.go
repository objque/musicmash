package db

import (
	"time"
)

type InternalRelease struct {
	ID         uint64    `json:"id"`
	ArtistID   int64     `json:"artist_id"`
	ArtistName string    `json:"artist_name"`
	Released   time.Time `json:"released"`
	Poster     string    `json:"poster"`
	Title      string    `json:"title"`
	ItunesID   string    `json:"itunes_id"`
	SpotifyID  string    `json:"spotify_id"`
	DeezerID   string    `json:"deezer_id"`
	Type       string    `json:"type"`
}

type InternalReleaseMgr interface {
	GetArtistInternalReleases(id int64) ([]*InternalRelease, error)
	GetUserInternalReleases(userName string, since, till *time.Time) ([]*InternalRelease, error)
}

func (mgr *AppDatabaseMgr) GetArtistInternalReleases(id int64) ([]*InternalRelease, error) {
	const query = `
SELECT releases.id,
       releases.artist_id,
       artists.name AS artist_name,
	   releases.released,
	   releases.poster,
	   releases.title,
	   releases.type,
	   itunes.store_id  AS itunes_id,
	   spotify.store_id AS spotify_id,
	   deezer.store_id  AS deezer_id
FROM releases AS releases
LEFT JOIN artists ON (
   releases.artist_id = artists.id
)
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
WHERE releases.artist_id = ?
GROUP BY releases.title
ORDER BY releases.released DESC
`
	releases := []*InternalRelease{}
	err := mgr.db.Raw(query, id).Scan(&releases).Error
	if err != nil {
		return nil, err
	}
	return releases, nil
}

func (mgr *AppDatabaseMgr) GetUserInternalReleases(userName string, since, till *time.Time) ([]*InternalRelease, error) {
	const query = `
SELECT releases.id,
       releases.artist_id,
       artists.name AS artist_name,
	   releases.released,
	   releases.poster,
	   releases.title,
	   releases.type,
	   itunes.store_id  AS itunes_id,
	   spotify.store_id AS spotify_id,
	   deezer.store_id  AS deezer_id
FROM releases AS releases
LEFT JOIN artists ON (
   releases.artist_id = artists.id
)
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
WHERE releases.artist_id IN (
   SELECT artist_id FROM subscriptions
   WHERE user_name = ?
) AND (
   releases.released >= ? AND releases.released < ?
)
GROUP BY releases.artist_id, releases.title
ORDER BY releases.released DESC
`
	releases := []*InternalRelease{}
	err := mgr.db.Raw(query, userName, since.Format("2006-01-02"), till.Format("2006-01-02")).Scan(&releases).Error
	if err != nil {
		return nil, err
	}
	return releases, nil
}
