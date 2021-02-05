package db

import (
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type InternalRelease struct {
	Released    time.Time `db:"released"`
	ReleaseDate string    `db:"-"`
	ArtistName  string    `db:"artist_name"`
	Poster      string    `db:"poster"`
	Title       string    `db:"title"`
	SpotifyID   string    `db:"spotify_id"`
	Type        string    `db:"type"`
	DurationMs  int64     `db:"duration_ms"`
	ID          uint64    `db:"id"`
	ArtistID    int64     `db:"artist_id"`
	TracksCount int32     `db:"tracks_count"`
	IsExplicit  bool      `db:"is_explicit"`
}

func (r *InternalRelease) MarshalJSON() ([]byte, error) {
	// TODO (m.kalinin): extract custom marshaler into repository package
	var release = struct {
		Released    string `json:"released"`
		ArtistName  string `json:"artist_name"`
		Poster      string `json:"poster"`
		Title       string `json:"title"`
		SpotifyID   string `json:"spotify_id"`
		Type        string `json:"type"`
		DurationMs  int64  `json:"duration_ms"`
		ID          uint64 `json:"id"`
		ArtistID    int64  `json:"artist_id"`
		TracksCount int32  `json:"tracks_count"`
		IsExplicit  bool   `json:"explicit"`
	}{
		Released:    r.Released.Format("2006-01-02"),
		ArtistName:  r.ArtistName,
		Poster:      r.Poster,
		Title:       r.Title,
		SpotifyID:   r.SpotifyID,
		Type:        r.Type,
		DurationMs:  r.DurationMs,
		ID:          r.ID,
		ArtistID:    r.ArtistID,
		TracksCount: r.TracksCount,
		IsExplicit:  r.IsExplicit,
	}

	return json.Marshal(&release)
}

type GetInternalReleaseOpts struct {
	Before      *uint64
	Limit       *uint64
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
		"releases.id",
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

//nolint:gocognit
func applyInternalReleasesFilters(query sq.SelectBuilder, opts *GetInternalReleaseOpts) sq.SelectBuilder {
	// we should choose only one filter for artists: artist_id or user subscriptions
	if opts.ArtistID != nil {
		query = query.Where("releases.artist_id = ?", *opts.ArtistID)
	}

	if opts.UserName != "" {
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

	if opts.Before != nil {
		query = query.Where("releases.id < ?", *opts.Before)
	}

	if opts.SortType != "" {
		// OrderByClause method generates incorrect query and we can't pass ASC/DESC as an arg
		query = query.OrderBy(
			fmt.Sprintf("releases.released %v", opts.SortType),
			fmt.Sprintf("releases.id %v", opts.SortType),
		)
	}

	if opts.Limit != nil {
		query = query.Limit(*opts.Limit)
	}

	return query
}
