package db

import (
	"database/sql"
	"time"
)

type Release struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	ArtistID  int64     `db:"artist_id"`
	Title     string    `db:"title"`
	Poster    string    `db:"poster"`
	Released  time.Time `db:"released"`
	StoreName string    `db:"store_name"`
	StoreID   string    `db:"store_id"`
	Type      string    `db:"type"`
	Explicit  bool      `db:"explicit"`
}

type ReleaseMgr interface {
	EnsureReleaseExists(release *Release) error
	CreateRelease(release *Release) error
	GetAllReleases() ([]*Release, error)
	FindReleases(condition map[string]interface{}) ([]*Release, error)
	InsertBatchNewReleases(releases []*Release) error
}

func (mgr *AppDatabaseMgr) EnsureReleaseExists(release *Release) error {
	const query = "select * from releases where store_id = $1 and store_name = $2 limit 1"

	res := Release{}
	err := mgr.newdb.Get(&res, query, release.StoreID, release.StoreName)
	if err == sql.ErrNoRows {
		return mgr.CreateRelease(release)
	}
	return err
}

func (mgr *AppDatabaseMgr) CreateRelease(release *Release) error {
	const query = "insert into releases (created_at, artist_id, title, poster, released, store_name, store_id, type, explicit) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	r, err := mgr.newdb.Exec(query, release.CreatedAt, release.ArtistID, release.Title, release.Poster,
		release.Released, release.StoreName, release.StoreID, release.Type, release.Explicit)
	if err != nil {
		return err
	}

	id, _ := r.LastInsertId()
	release.ID = uint64(id)
	return nil
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
	const query = "select * from releases where artist_id = $1 and title = ?"

	releases := []*Release{}
	err := mgr.newdb.Select(&releases, query, artistID, title)
	if err != nil {
		return nil, err
	}

	return releases, nil
}
