package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/musicmash/musicmash/internal/api/httputils"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

const (
	PathArtists = "/artists"
)

type ArtistsController struct{}

func NewArtistsController() *ArtistsController {
	return &ArtistsController{}
}

func (c *ArtistsController) Register(router chi.Router) {
	router.Route(PathArtists, func(r chi.Router) {
		r.Post("/", c.addArtist)
		r.Get("/{id}", c.getArtist)
	})
}

func (c *ArtistsController) addArtist(w http.ResponseWriter, r *http.Request) {
	artist := db.Artist{}
	err := json.NewDecoder(r.Body).Decode(&artist)
	if err != nil {
		httputils.WriteError(w, errors.New("invalid body"))
		return
	}

	// TODO (m.kalinin): when project was moved from gorm to sqlx we have broke method db.CreateArtist
	// it requires id for insert an new entity. so for temporary reason code below was commented.
	// do not allow override ID
	//if artist.ID != 0 {
	//	httputils.WriteError(w, errors.New("artist id should be empty"))
	//	return
	//}

	if artist.Name == "" {
		httputils.WriteError(w, errors.New("artist name didn't provided"))
		return
	}

	err = db.Mgr.EnsureArtistExists(&artist)
	if err != nil {
		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusCreated, &artist)
}

func (c *ArtistsController) getArtist(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputils.WriteError(w, errors.New("wrong id"))
		return
	}

	artist, err := db.Mgr.GetArtist(id)
	if err != nil {
		if err == sql.ErrNoRows {
			httputils.WriteError(w, errors.New("artist not found"))
			return
		}

		httputils.WriteInternalError(w)
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &artist)
}
