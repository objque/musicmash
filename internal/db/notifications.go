package db

import "time"

type Notification struct {
	ID        int
	Date      time.Time
	UserName  string
	ReleaseID uint64
	IsComing  bool
}

type NotificationMgr interface {
	CreateNotification(notification *Notification) error
	GetNotificationsForUser(userName string) ([]*Notification, error)
}

func (mgr *AppDatabaseMgr) GetNotificationsForUser(userName string) ([]*Notification, error) {
	notifications := []*Notification{}
	if err := mgr.db.Where("user_name = ?", userName).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (mgr *AppDatabaseMgr) CreateNotification(notification *Notification) error {
	return mgr.db.Create(&notification).Error
}
