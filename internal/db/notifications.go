package db

import (
	"time"

	"github.com/musicmash/musicmash/internal/log"
	"github.com/pkg/errors"
)

type Notification struct {
	ID        int       `gorm:"primary_key" sql:"AUTO_INCREMENT"`
	Date      time.Time `gorm:"unique_index:idx_notify_date_user_release"`
	UserName  string    `gorm:"unique_index:idx_notify_date_user_release"`
	ReleaseID uint64    `gorm:"unique_index:idx_notify_date_user_release"`
}

type NotificationMgr interface {
	CreateNotification(notification *Notification) error
	GetNotificationsForUser(userName string) ([]*Notification, error)
	MarkReleasesAsDelivered(userName string, releases []*Release)
	IsUserAlreadyNotified(userName string, release *Release) (bool, error)
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

func (mgr *AppDatabaseMgr) IsUserAlreadyNotified(userName string, release *Release) (bool, error) {
	count := 0
	query := mgr.db.Table("notifications")
	const whereCondition = "user_name = ? and release_id = ? and date = ?"
	err := query.Where(whereCondition, userName, release.ID, time.Now().UTC().Truncate(time.Hour*24)).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (mgr *AppDatabaseMgr) MarkReleasesAsDelivered(userName string, releases []*Release) {
	for _, release := range releases {
		notified, err := mgr.IsUserAlreadyNotified(userName, release)
		if notified {
			continue
		}
		if err != nil {
			log.Error(err)
			continue
		}

		notification := Notification{
			UserName:  userName,
			ReleaseID: release.ID,
			Date:      time.Now().UTC().Truncate(time.Hour * 24),
		}
		if err := mgr.CreateNotification(&notification); err != nil {
			log.Error(errors.Wrapf(err, "tried to save notification for user '%s' about release_id '%v'", userName, release.ID))
		}
	}
}
