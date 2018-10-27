package validators

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

func IsUserExits(w http.ResponseWriter, name string) error {
	_, err := db.DbMgr.FindUserByName(name)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			w.WriteHeader(http.StatusNotFound)
			return err
		}

		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return err
}
