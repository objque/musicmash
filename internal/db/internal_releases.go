package db

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type InternalRelease struct {
	Released    time.Time `json:"released"     db:"released"`
	ArtistName  string    `json:"artist_name"  db:"artist_name"`
	Poster      string    `json:"poster"       db:"poster"`
	Title       string    `json:"title"        db:"title"`
	SpotifyID   string    `json:"spotify_id"   db:"spotify_id"`
	Type        string    `json:"type"         db:"type"`
	DurationMs  int64     `json:"duration_ms"  db:"duration_ms"`
	ArtistID    int64     `json:"artist_id"    db:"artist_id"`
	TracksCount int32     `json:"tracks_count" db:"tracks_count"`
	IsExplicit  bool      `json:"explicit"     db:"is_explicit"`
}

type GetInternalReleaseOpts struct {
	Limit       *uint64
	Offset      *uint64
	ArtistID    *int64
	UserName    string
	ReleaseType string
	SortType    string
	Title       string
	Explicit    *bool
	Since       *time.Time
	Till        *time.Time
}

func (mgr *AppDatabaseMgr) GetInternalReleases(opts *GetInternalReleaseOpts) ([]*InternalRelease, error) {
	query := sq.Select(
		"releases.artist_id",
		"artists.name AS artist_name",
		"releases.released",
		"releases.poster",
		"releases.title",
		"releases.type",
		"releases.is_explicit",
		"releases.tracks_count",
		"releases.duration_ms",
		"releases.spotify_id").
		From("releases").
		LeftJoin("artists ON (releases.artist_id = artists.id)")

	if opts != nil {
		query = applyInternalReleasesFilters(query, opts)
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	releases := make([]*InternalRelease, 0)
	if err := mgr.newdb.Select(&releases, sql, args...); err != nil {
		return nil, err
	}

	return releases, nil
}

func applyInternalReleasesFilters(query sq.SelectBuilder, opts *GetInternalReleaseOpts) sq.SelectBuilder {
	// we should choose only one filter for artists: artist_id or user subscriptions
	if opts.ArtistID != nil {
		query = query.Where("releases.artist_id = ?", *opts.ArtistID)
	} else if opts.UserName != "" {
		const format = "SELECT artist_id FROM subscriptions WHERE user_name = '%v'"
		subQ := fmt.Sprintf(format, opts.UserName)
		query = query.Where(fmt.Sprintf("releases.artist_id IN (%v)", subQ))
	}

	if opts.ReleaseType != "" {
		query = query.Where("releases.type = ?", opts.ReleaseType)
	}

	if opts.Title != "" {
		query = query.Where("releases.title like ?", fmt.Sprint("%", opts.Title, "%"))
	}

	if opts.Since != nil {
		query = query.Where("releases.released >= ?", opts.Since.Format("2006-01-02"))
	}

	if opts.Till != nil {
		query = query.Where("releases.released < ?", opts.Till.Format("2006-01-02"))
	}

	if opts.Explicit != nil {
		query = query.Where("releases.is_explicit = ?", *opts.Explicit)
	}

	if opts.SortType != "" {
		// OrderByClause method generates incorrect query and we can't pass ASC/DESC as an arg
		query = query.OrderBy(fmt.Sprintf("releases.released %v", opts.SortType))
	}

	if opts.Offset != nil {
		query = query.Offset(*opts.Offset)
	}

	if opts.Limit != nil {
		query = query.Limit(*opts.Limit)
	}

	return query
}
