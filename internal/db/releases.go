package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Release struct {
	ID         uint64
	CreatedAt  time.Time
	ArtistName string
	Title      string
	Poster     string
	Released   time.Time
	StoreName  string `gorm:"unique_index:idx_rel_store_name_store_id"`
	StoreID    string `gorm:"unique_index:idx_rel_store_name_store_id"`
}

type ReleaseMgr interface {
	EnsureReleaseExists(release *Release) error
	GetAllReleases() ([]*Release, error)
}

func (mgr *AppDatabaseMgr) EnsureReleaseExists(release *Release) error {
	res := Release{}
	err := mgr.db.Where("store_id = ? and store_name = ?", release.StoreID, release.StoreName).First(&res).Error
	if gorm.IsRecordNotFoundError(err) {
		return mgr.db.Create(release).Error
	}
	return err
}

func (mgr *AppDatabaseMgr) GetAllReleases() ([]*Release, error) {
	var releases = make([]*Release, 0)
	return releases, mgr.db.Find(&releases).Error
}
