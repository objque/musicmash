package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	body := CreateUserScheme{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(body.UserName)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := db.DbMgr.FindUserByName(body.UserName)
	// already exists
	if err == nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// err while processing query
	if err != gorm.ErrRecordNotFound {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := db.DbMgr.CreateUser(&db.User{Name: body.UserName}); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "user_name")
	_, err := db.DbMgr.FindUserByName(userName)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	if gorm.IsRecordNotFoundError(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
}
