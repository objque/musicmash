package db

import (
	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/log"
)

type StoreType struct {
	Name string `gorm:"primary_key"`
}

type StoreTypeMgr interface {
	IsStoreExists(name string) bool
	EnsureStoreExists(name string) error
}

func (mgr *AppDatabaseMgr) IsStoreExists(name string) bool {
	store := StoreType{}
	if err := mgr.db.Where("name = ?", name).First(&store).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}

		log.Error(err)
		return false
	}
	return true
}

func (mgr *AppDatabaseMgr) EnsureStoreExists(name string) error {
	if !mgr.IsStoreExists(name) {
		return mgr.db.Create(&StoreType{Name: name}).Error
	}
	return nil
}
