package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/musicmash/musicmash/internal/api/httputils"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

const (
	PathReleases = "/releases"
)

type ReleasesController struct{}

func NewReleasesController() *ReleasesController {
	return &ReleasesController{}
}

func (rc *ReleasesController) Register(router chi.Router) {
	router.Route(PathReleases, func(r chi.Router) {
		r.Get("/", rc.getReleases)
	})
}

//nolint:gocyclo,gocognit
func (rc *ReleasesController) getReleases(w http.ResponseWriter, r *http.Request) {
	var defaultMaxLimit uint64 = 100
	opts := db.GetInternalReleaseOpts{
		Limit:    &defaultMaxLimit,
		SortType: "DESC",
	}

	opts.UserName, _ = GetUser(r)

	// todo: extract all query parsers
	if v := r.URL.Query().Get("before"); v != "" {
		before, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			httputils.WriteError(w, errors.New("before must be int and greater than 0"))
			return
		}

		opts.Before = &before
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		offset, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			httputils.WriteError(w, errors.New("offset must be int and greater than 0"))
			return
		}

		opts.Offset = &offset
	}

	if v := r.URL.Query().Get("limit"); v != "" {
		limit, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			httputils.WriteError(w, errors.New("limit must be int and greater than 0, but less than 100"))
			return
		}

		if limit > defaultMaxLimit {
			httputils.WriteError(w, errors.New("limit must be int and greater than 0, but less than 100"))
			return
		}

		opts.Limit = &limit
	}

	if v := r.URL.Query().Get("artist_id"); v != "" {
		artistID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			httputils.WriteError(w, errors.New("artist_id must be int and greater than 0"))
			return
		}

		opts.ArtistID = &artistID
	}

	if v := r.URL.Query().Get("explicit"); v != "" {
		explicit, err := strconv.ParseBool(v)
		if err != nil {
			httputils.WriteError(w, errors.New("explicit must be true or false"))
			return
		}

		opts.Explicit = &explicit
	}

	if opts.ReleaseType = r.URL.Query().Get("type"); opts.ReleaseType != "" {
		if opts.ReleaseType != "album" && opts.ReleaseType != "song" && opts.ReleaseType != "music-video" {
			httputils.WriteError(w, errors.New("type must be one of {album,song,music-video}"))
			return
		}
	}

	if v := strings.ToUpper(r.URL.Query().Get("sort_type")); v != "" {
		if v != "ASC" && v != "DESC" {
			httputils.WriteError(w, errors.New("sort_type must be one of {asc,desc}"))
			return
		}

		opts.SortType = v
	}

	if v := r.URL.Query().Get("since"); v != "" {
		since, err := time.Parse("2006-01-02", v)
		if err != nil {
			httputils.WriteError(w, errors.New("since must be in format YYYY-MM-DD"))
			return
		}

		opts.Since = &since
	}

	if v := r.URL.Query().Get("till"); v != "" {
		till, err := time.Parse("2006-01-02", v)
		if err != nil {
			httputils.WriteError(w, errors.New("till must be in format YYYY-MM-DD"))
			return
		}

		opts.Till = &till
	}

	if opts.Since != nil && opts.Till != nil && opts.Since.After(*opts.Till) {
		httputils.WriteError(w, errors.New("since must be before till"))
		return
	}

	releases, err := db.Mgr.GetInternalReleases(&opts)
	if err != nil {
		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &releases)
}
