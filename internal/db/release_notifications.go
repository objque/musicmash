package db

import (
	"time"

	sq "github.com/Masterminds/squirrel"
)

type ReleaseNotification struct {
	ArtistID   int64     `db:"artist_id"`
	ArtistName string    `db:"artist_name"`
	CreatedAt  time.Time `db:"created_at"`
	Released   time.Time `db:"released"`
	Poster     string    `db:"poster"`
	Title      string    `db:"title"`
	UserName   string    `db:"user_name"`
	SpotifyID  string    `db:"spotify_id"`
	Type       string    `db:"type"`
	IsExplicit bool      `db:"is_explicit"`
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
		"releases.is_explicit",
		"releases.spotify_id").
		From("releases AS releases").
		JoinClause(`INNER JOIN subscriptions ON (
			subscriptions.artist_id = releases.artist_id
		)`).
		LeftJoin("artists ON (releases.artist_id = artists.id)").
		Where("releases.created_at >= ?", since.Format("2006-01-02T15:04:05")).
		GroupBy(
			"subscriptions.user_name",
			"releases.artist_id",
			"artist_name",
			"releases.released",
			"releases.poster",
			"releases.title",
			"releases.type",
			"releases.is_explicit",
			"spotify_id").
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
