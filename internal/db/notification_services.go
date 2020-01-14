package db

import (
	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/log"
)

type NotificationService struct {
	ID string
}

type NotificationServiceMgr interface {
	IsNotificationServiceExists(id string) bool
	EnsureNotificationServiceExists(id string) error
}

func (mgr *AppDatabaseMgr) IsNotificationServiceExists(id string) bool {
	service := NotificationService{}
	if err := mgr.db.Where("id = ?", id).First(&service).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}

		log.Error(err)
		return false
	}
	return true
}

func (mgr *AppDatabaseMgr) EnsureNotificationServiceExists(id string) error {
	if !mgr.IsNotificationServiceExists(id) {
		return mgr.db.Create(&NotificationService{ID: id}).Error
	}
	return nil
}
