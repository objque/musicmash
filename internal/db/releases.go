package db

import (
	"database/sql"
	"time"
)

type Release struct {
	ID          uint64         `db:"id"`
	CreatedAt   time.Time      `db:"created_at"`
	ArtistID    int64          `db:"artist_id"`
	Type        string         `db:"type"`
	TracksCount int32          `db:"tracks_count"`
	DurationMs  int64          `db:"duration_ms"`
	Title       string         `db:"title"`
	Released    time.Time      `db:"released"`
	Explicit    bool           `db:"is_explicit"`
	Popularity  sql.NullInt32  `db:"popularity"`
	Poster      sql.NullString `db:"poster"`
	SpotifyID   string         `db:"spotify_id"`
}

func (mgr *AppDatabaseMgr) EnsureReleaseExists(release *Release) error {
	const query = "select * from releases where spotify_id = $1 limit 1"

	res := Release{}
	err := mgr.newdb.Get(&res, query, release.SpotifyID)
	if err == sql.ErrNoRows {
		return mgr.CreateRelease(release)
	}
	return err
}

func (mgr *AppDatabaseMgr) CreateRelease(release *Release) error {
	const query = "insert into releases (created_at, artist_id, title, poster, released, spotify_id, type, is_explicit, tracks_count, duration_ms) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id"

	return mgr.newdb.QueryRow(query, release.CreatedAt, release.ArtistID, release.Title, release.Poster,
		release.Released, release.SpotifyID, release.Type, release.Explicit, release.TracksCount, release.DurationMs).Scan(&release.ID)
}

func (mgr *AppDatabaseMgr) GetAllReleases() ([]*Release, error) {
	const query = "select * from releases"

	var releases = []*Release{}
	err := mgr.newdb.Select(&releases, query)
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func (mgr *AppDatabaseMgr) FindReleases(artistID int64, title string) ([]*Release, error) {
	const query = "select * from releases where artist_id = $1 and title = $2"

	releases := []*Release{}
	err := mgr.newdb.Select(&releases, query, artistID, title)
	if err != nil {
		return nil, err
	}

	return releases, nil
}
