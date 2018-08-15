package validators

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
)

func IsUserExits(w http.ResponseWriter, userID string) error {
	_, err := db.DbMgr.FindUserByID(userID)
	if err != nil {
		statusCode := http.StatusNotFound
		if !gorm.IsRecordNotFoundError(err) {
			statusCode = http.StatusInternalServerError
		}

		log.Error(err)
		w.WriteHeader(statusCode)
	}
	return err
}
