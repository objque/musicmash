package db

import "time"

type Notification struct {
	ID        int       `db:"id"`
	Date      time.Time `db:"date"`
	UserName  string    `db:"user_name"`
	ReleaseID uint64    `db:"release_id"`
	IsComing  bool      `db:"is_coming"`
}

func (mgr *AppDatabaseMgr) GetNotificationsForUser(userName string) ([]*Notification, error) {
	const query = "select * from notifications where user_name = $1"

	notifications := []*Notification{}
	err := mgr.newdb.Select(&notifications, query, userName)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (mgr *AppDatabaseMgr) CreateNotification(notification *Notification) error {
	const query = "insert into notifications (date, user_name, release_id, is_coming) values ($1, $2, $3, $4)"

	_, err := mgr.newdb.Exec(query, notification.Date, notification.UserName, notification.ReleaseID, notification.IsComing)

	return err
}
