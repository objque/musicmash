package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
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

func (rc *ReleasesController) getReleases(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("artist_id") != "" {
		rc.getReleasesByArtist(w, r)
		return
	}

	rc.getReleasesForUser(w, r)
}

func (rc *ReleasesController) getReleasesByArtist(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("artist_id"), 10, 64)
	if err != nil {
		httputils.WriteError(w, errors.New("wrong id"))
		return
	}

	_, err = db.DbMgr.GetArtist(id)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			httputils.WriteError(w, errors.New("artist not found"))
			return
		}

		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	releases, err := db.DbMgr.GetArtistInternalReleases(id)
	if err != nil {
		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &releases)
}

func (rc *ReleasesController) getReleasesForUser(w http.ResponseWriter, r *http.Request) {
	userName, err := GetUser(r)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	since, _ := httputils.GetQueryTimeWithLayout(r, "since", "2006-01-02")
	if since == nil {
		s := now.New(time.Now()).BeginningOfWeek()
		since = &s
	}

	till, _ := httputils.GetQueryTimeWithLayout(r, "till", "2006-01-02")
	if till == nil {
		t := since.AddDate(0, 0, 7)
		till = &t
	}

	if since.After(*till) {
		httputils.WriteError(w, errors.New("since must be before till"))
		return
	}

	releases, err := db.DbMgr.GetUserInternalReleases(userName, since, till)
	if err != nil {
		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &releases)
}
