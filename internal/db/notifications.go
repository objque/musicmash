package db

import (
	"time"

	"github.com/musicmash/musicmash/internal/log"
	"github.com/pkg/errors"
)

type Notification struct {
	ID        int `gorm:"primary_key" sql:"AUTO_INCREMENT"`
	CreatedAt time.Time
	UserName  string `gorm:"unique_index:idx_notify_user_release"`
	ReleaseID uint64 `gorm:"unique_index:idx_notify_user_release"`
}

type NotificationMgr interface {
	GetNotificationsForUser(userName string) ([]*Notification, error)
	MarkReleasesAsDelivered(userName string, releases []*Release)
}

func (mgr *AppDatabaseMgr) GetNotificationsForUser(userName string) ([]*Notification, error) {
	notifications := []*Notification{}
	if err := mgr.db.Where("user_name = ?", userName).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (mgr *AppDatabaseMgr) MarkReleasesAsDelivered(userName string, releases []*Release) {
	for _, release := range releases {
		if err := mgr.db.Create(&Notification{UserName: userName, ReleaseID: release.ID}).Error; err != nil {
			log.Error(errors.Wrapf(err, "tried to save notification for user '%s' about release_id '%v'", userName, release.ID))
		}
	}
}
