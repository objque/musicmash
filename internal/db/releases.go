package db

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/log"
)

type Release struct {
	ID         int64     `gorm:"primary_key" sql:"AUTO_INCREMENT" json:"-"`
	CreatedAt  time.Time `json:"-"`
	Date       time.Time `gorm:"not null" sql:"index" json:"date"`
	ArtistName string    `json:"artist_name"`
	StoreID    uint64    `sql:"index" gorm:"index:idx_store_id" json:"store_id"`
	Stores     []Store   `json:"stores"`
}

type ReleaseMgr interface {
	CreateRelease(release *Release) error
	FindRelease(artist string, id uint64) (*Release, error)
	IsReleaseExists(id uint64) bool
	GetAllReleases() ([]*Release, error)
	EnsureReleaseExists(release *Release) error
	GetReleasesForUserFilterByPeriod(userID string, since, till time.Time) ([]*Release, error)
	GetReleasesForUserSince(userID string, since time.Time) ([]*Release, error)
}

func (mgr *AppDatabaseMgr) FindRelease(artist string, id uint64) (*Release, error) {
	release := Release{}
	if err := mgr.db.Where("artist_name = ? and store_id = ?", artist, id).First(&release).Error; err != nil {
		return nil, err
	}
	return &release, nil
}

func (mgr *AppDatabaseMgr) IsReleaseExists(id uint64) bool {
	release := Release{}
	if err := mgr.db.Where("store_id = ?", id).First(&release).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}

		log.Error(err)
		return false
	}
	return true
}

func (mgr *AppDatabaseMgr) GetAllReleases() ([]*Release, error) {
	var releases = make([]*Release, 0)
	return releases, mgr.db.Find(&releases).Error
}

func (mgr *AppDatabaseMgr) CreateRelease(release *Release) error {
	return mgr.db.Create(release).Error
}

func (mgr *AppDatabaseMgr) EnsureReleaseExists(release *Release) error {
	if !mgr.IsReleaseExists(release.StoreID) {
		return mgr.CreateRelease(release)
	}
	return nil
}

func (mgr *AppDatabaseMgr) GetReleasesForUserFilterByPeriod(userID string, since, till time.Time) ([]*Release, error) {
	// select * from releases where artist_name in (filter by user_subs) and filter by since/till dates;
	// select * from stores where release_id in (ids from query above);
	releases := []*Release{}
	const query = "select artist_name from subscriptions where user_id = ?"
	innerQuery := mgr.db.Raw(query, userID).QueryExpr()
	where := mgr.db.Where("artist_name in (?) and date >= ? and date <= ?", innerQuery, since, till)
	if err := where.Preload("Stores").Find(&releases).Error; err != nil {
		return nil, err
	}
	return releases, nil
}

func (mgr *AppDatabaseMgr) GetReleasesForUserSince(userID string, since time.Time) ([]*Release, error) {
	// select * from releases where artist_name in (filter by user_subs) and filter by since date;
	// select * from stores where release_id in (ids from query above);
	releases := []*Release{}
	const query = "select artist_name from subscriptions where user_id = ?"
	innerQuery := mgr.db.Raw(query, userID).QueryExpr()
	where := mgr.db.Where("artist_name in (?) and date >= ?", innerQuery, since)
	if err := where.Preload("Stores").Find(&releases).Error; err != nil {
		return nil, err
	}
	return releases, nil
}
