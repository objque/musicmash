package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
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
		r.Post("/associate", c.associateArtist)
		r.Get("/", c.searchArtist)
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

	// do not allow override ID
	if artist.Name == "" {
		httputils.WriteError(w, errors.New("artist name didn't provided"))
		return
	}
	artist.ID = 0

	err = db.DbMgr.EnsureArtistExists(&artist)
	if err != nil {
		httputils.WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusCreated, &artist)
}

func (c *ArtistsController) associateArtist(w http.ResponseWriter, r *http.Request) {
	info := db.Association{}
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		httputils.WriteError(w, errors.New("invalid body"))
		return
	}

	if exist := db.DbMgr.IsStoreExists(info.StoreName); !exist {
		httputils.WriteError(w, errors.New("store not found"))
		return
	}

	if exist := db.DbMgr.IsAssociationExists(info.StoreName, info.StoreID); exist {
		httputils.WriteError(w, errors.New("artist already associated"))
		return
	}

	_, err = db.DbMgr.GetArtist(info.ArtistID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			httputils.WriteError(w, errors.New("artist not found"))
			return
		}

		httputils.WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	err = db.DbMgr.EnsureAssociationExists(info.ArtistID, info.StoreName, info.StoreID)
	if err != nil {
		httputils.WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusCreated, &info)
}

func (c *ArtistsController) searchArtist(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		httputils.WriteError(w, errors.New("query argument name didn't provided"))
		return
	}

	artists, err := db.DbMgr.SearchArtists(name)
	if err != nil {
		httputils.WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &artists)
}

func (c *ArtistsController) getArtist(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputils.WriteError(w, errors.New("wrong id"))
		return
	}

	artist, err := db.DbMgr.GetArtist(id)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			httputils.WriteError(w, errors.New("artist not found"))
			return
		}

		httputils.WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
		log.Error(err)
		return
	}

	includeAlbums, _ := strconv.ParseBool(r.URL.Query().Get("withAlbums"))
	if includeAlbums {
		artist.Albums, err = db.DbMgr.GetAlbums(id)
		if err != nil {
			httputils.WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal"))
			log.Error(err)
			return
		}
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &artist)
}
