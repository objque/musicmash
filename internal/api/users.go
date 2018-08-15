package api

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	body := CreateUserScheme{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := db.DbMgr.FindUserByID(body.UserID)
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

	if err := db.DbMgr.CreateUser(&db.User{ID: body.UserID}); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
