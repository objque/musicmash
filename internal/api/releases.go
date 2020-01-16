package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
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
		r.Get("/", rc.getReleasesForUser)
	})
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
