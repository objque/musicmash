package db

import (
	"time"

	sq "github.com/Masterminds/squirrel"
)

type ReleaseNotification struct {
	ArtistID   int64     `json:"artist_id"   db:"artist_id"`
	ArtistName string    `json:"artist_name" db:"artist_name"`
	CreatedAt  time.Time `json:"-"           db:"created_at"`
	Released   time.Time `json:"released"    db:"released"`
	Poster     string    `json:"poster"      db:"poster"`
	Title      string    `json:"title"       db:"title"`
	UserName   string    `json:"user_name"   db:"user_name"`
	ItunesID   *string   `json:"itunes_id"   db:"itunes_id"`
	SpotifyID  *string   `json:"spotify_id"  db:"spotify_id"`
	DeezerID   *string   `json:"deezer_id"   db:"deezer_id"`
	Type       string    `json:"type"        db:"type"`
	Explicit   bool      `json:"explicit"    db:"explicit"`
}

func (mgr *AppDatabaseMgr) GetReleaseNotifications(since time.Time) ([]*ReleaseNotification, error) {
	query := sq.Select(
		"subscriptions.user_name",
		"releases.artist_id",
		"artists.name AS artist_name",
		"releases.released",
		"releases.poster",
		"releases.title",
		"releases.type",
		"releases.explicit",
		"itunes.store_id AS itunes_id",
		"spotify.store_id AS spotify_id",
		"deezer.store_id AS deezer_id").
		From("releases AS releases").
		JoinClause(`INNER JOIN subscriptions ON (
			subscriptions.artist_id = releases.artist_id
		)`).
		LeftJoin("artists ON (releases.artist_id = artists.id)").
		LeftJoin(`releases AS itunes ON (
			releases.artist_id = itunes.artist_id
			AND releases.title = itunes.title
			AND itunes.store_name = 'itunes'
		)`).
		LeftJoin(`releases AS spotify ON (
			releases.artist_id = spotify.artist_id
			AND releases.title = spotify.title
			AND spotify.store_name = 'spotify'
		)`).
		LeftJoin(`releases AS deezer ON (
			releases.artist_id = deezer.artist_id
			AND releases.title = deezer.title
			AND deezer.store_name = 'deezer'
		)`).
		Where("releases.created_at >= ?", since.Format("2006-01-02")).
		GroupBy(
			"subscriptions.user_name",
			"releases.artist_id",
			"artist_name",
			"releases.released",
			"releases.poster",
			"releases.title",
			"releases.type",
			"releases.explicit",
			"itunes_id",
			"spotify_id",
			"deezer_id").
		OrderBy("user_name, released ASC")

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	releases := make([]*ReleaseNotification, 0)
	if err := mgr.newdb.Select(&releases, sql, args...); err != nil {
		return nil, err
	}

	return releases, nil
}
