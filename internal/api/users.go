package api

import (
	"encoding/json"
	"net/http"

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

	_, err := db.DbMgr.FindUserByName(body.UserName)
	// already exists
	if err == nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// err while processing query
	if err != nil && err != gorm.ErrRecordNotFound {
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
