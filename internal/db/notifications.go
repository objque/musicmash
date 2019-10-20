package db

import "time"

type Notification struct {
	ID        int `gorm:"primary_key" sql:"AUTO_INCREMENT"`
	Date      time.Time
	UserName  string `gorm:"unique_index:idx_user_name_release_id_is_coming"`
	ReleaseID uint64 `gorm:"unique_index:idx_user_name_release_id_is_coming"`
	IsComing  bool   `gorm:"unique_index:idx_user_name_release_id_is_coming" gorm:"default:0'"`
}

type NotificationMgr interface {
	CreateNotification(notification *Notification) error
	GetNotificationsForUser(userName string) ([]*Notification, error)
	IsUserNotified(userName string, releaseID uint64, isComing bool) (*Notification, error)
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

func (mgr *AppDatabaseMgr) IsUserNotified(userName string, releaseID uint64, isComing bool) (*Notification, error) {
	notification := Notification{}
	const query = "user_name = ? and release_id = ? and is_coming = ?"
	err := mgr.db.Where(query, userName, releaseID, isComing).First(&notification).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}
