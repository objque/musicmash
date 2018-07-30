package db

import (
	"time"
)

type Release struct {
	ID         int64 `gorm:"primary_key" sql:"AUTO_INCREMENT"`
	CreatedAt  time.Time
	ArtistName string
	StoreURL   string
}

type ReleaseMgr interface {
	CreateRelease(release *Release) error
	FindRelease(artist string, storeURL string) (*Release, error)
	GetAllReleases() ([]*Release, error)
	EnsureReleaseExists(release *Release) error
}

func (mgr *AppDatabaseMgr) FindRelease(artist string, storeURL string) (*Release, error) {
	release := Release{}
	if err := mgr.db.Where("artist_name = ? and store_url = ?", artist, storeURL).First(&release).Error; err != nil {
		return nil, err
	}
	return &release, mgr.db.Find(&release).Error
}

func (mgr *AppDatabaseMgr) GetAllReleases() ([]*Release, error) {
	var releases = make([]*Release, 0)
	return releases, mgr.db.Find(&releases).Error
}

func (mgr *AppDatabaseMgr) CreateRelease(release *Release) error {
	return mgr.db.Create(release).Error
}

func (mgr *AppDatabaseMgr) EnsureReleaseExists(release *Release) error {
	_, err := mgr.FindRelease(release.ArtistName, release.StoreURL)
	if err != nil {
		return mgr.CreateRelease(release)
	}
	return nil
}
