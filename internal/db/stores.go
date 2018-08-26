package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/log"
)

type Store struct {
	ID        uint64 `json:"-" gorm:"primary_key" sql:"AUTO_INCREMENT"`
	StoreID   string `json:"store_id"`
	StoreType string `json:"type" sql:"index" gorm:"unique_index:idx_store_type_release_id" json:"store_type"`
	ReleaseID int64  `json:"-" sql:"index" gorm:"unique_index:idx_store_type_release_id" json:"store_type"`
}

func (s *Store) GetName() string {
	switch s.StoreType {
	case "yandex":
		return "Yandex.Music"
	default:
		return "iTunes"
	}
}

func (s *Store) GetLink() string {
	switch s.StoreType {
	case "yandex":
		return fmt.Sprintf("https://music.yandex.ru/album/%s", s.StoreID)
	default:
		return fmt.Sprintf("https://itunes.apple.com/us/album/%s?uo=4", s.StoreID)
	}
}

type StoreMgr interface {
	IsReleaseExistsInStore(store string, storeID string) bool
	EnsureReleaseExistsInStore(store string, storeID string, releaseID int64) error
}

func (mgr *AppDatabaseMgr) IsReleaseExistsInStore(store string, storeID string) bool {
	s := Store{}
	if err := mgr.db.Where("store_type = ? and store_id = ?", store, storeID).First(&s).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}

		log.Error(err)
		return false
	}
	return true
}

func (mgr *AppDatabaseMgr) EnsureReleaseExistsInStore(store string, storeID string, releaseID int64) error {
	if !mgr.IsReleaseExistsInStore(store, storeID) {
		return mgr.db.Create(&Store{StoreType: store, StoreID: storeID, ReleaseID: releaseID}).Error
	}
	return nil
}
