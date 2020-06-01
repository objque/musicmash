package db

import (
	"database/sql"

	"github.com/musicmash/musicmash/internal/log"
)

type NotificationService struct {
	ID string
}

func (mgr *AppDatabaseMgr) IsNotificationServiceExists(id string) bool {
	const query = "select * from notification_services where id = $1"

	service := NotificationService{}
	err := mgr.newdb.Get(&service, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}

		log.Error(err)
		return false
	}
	return true
}

func (mgr *AppDatabaseMgr) EnsureNotificationServiceExists(id string) error {
	if mgr.IsNotificationServiceExists(id) {
		return nil
	}

	const query = "insert into notification_services (id) values ($1)"

	_, err := mgr.newdb.Exec(query, id)

	return err
}
