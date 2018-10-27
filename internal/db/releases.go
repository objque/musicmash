package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ReleaseStore struct {
	StoreName string `json:"name"`
	StoreID   string `json:"id"`
}

type Release struct {
	ID         uint64          `json:"-"`
	CreatedAt  time.Time       `json:"-"`
	ArtistName string          `json:"artist_name"`
	Title      string          `json:"title"`
	Poster     string          `json:"poster"`
	Released   time.Time       `gorm:"type:datetime" json:"released"`
	StoreName  string          `gorm:"unique_index:idx_rel_store_name_store_id" json:"-"`
	StoreID    string          `gorm:"unique_index:idx_rel_store_name_store_id" json:"-"`
	Stores     []*ReleaseStore `gorm:"-" json:"stores"`
}

type ReleaseMgr interface {
	EnsureReleaseExists(release *Release) error
	GetAllReleases() ([]*Release, error)
	GetReleasesForUserFilterByPeriod(userName string, since, till time.Time) ([]*Release, error)
	GetReleasesForUserSince(userName string, since time.Time) ([]*Release, error)
	FindReleases(condition map[string]interface{}) ([]*Release, error)
	FindNewReleases(date time.Time) ([]*Release, error)
	FindNewReleasesForUser(userName string, date time.Time) ([]*Release, error)
	UpdateRelease(release *Release) error
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

func (mgr *AppDatabaseMgr) GetReleasesForUserFilterByPeriod(userName string, since, till time.Time) ([]*Release, error) {
	// inner query: select artist_name from subscriptions where user_name = XXX
	// select * from releases where artist_name in (INNER) and and released >= ? and released <= ?
	releases := []*Release{}
	const query = "select artist_name from subscriptions where user_name = ?"
	innerQuery := mgr.db.Raw(query, userName).QueryExpr()
	where := mgr.db.Where("artist_name in (?) and released >= ? and released <= ?", innerQuery, since, till)
	if err := where.Find(&releases).Error; err != nil {
		return nil, err
	}
	return releases, nil
}

func (mgr *AppDatabaseMgr) GetReleasesForUserSince(userName string, since time.Time) ([]*Release, error) {
	// inner query: select artist_name from subscriptions where user_name = XXX
	// select * from releases where artist_name in (INNER) and and released >= ?
	releases := []*Release{}
	const query = "select artist_name from subscriptions where user_name = ?"
	innerQuery := mgr.db.Raw(query, userName).QueryExpr()
	where := mgr.db.Where("artist_name in (?) and released >= ?", innerQuery, since)
	if err := where.Find(&releases).Error; err != nil {
		return nil, err
	}
	return releases, nil
}

func (mgr *AppDatabaseMgr) FindNewReleases(date time.Time) ([]*Release, error) {
	releases := []*Release{}
	if err := mgr.db.Where("created_at >= ?", date).Find(&releases).Error; err != nil {
		return nil, err
	}
	return releases, nil
}

func (mgr *AppDatabaseMgr) FindNewReleasesForUser(userName string, date time.Time) ([]*Release, error) {
	// inner query: select artist_name from subscriptions where user_name = XXX
	// select * from releases where artist_name in (INNER) and and created_at >= ?
	releases := []*Release{}
	const query = "select artist_name from subscriptions where user_name = ?"
	innerQuery := mgr.db.Raw(query, userName).QueryExpr()
	where := mgr.db.Where("artist_name in (?) and created_at >= ?", innerQuery, date)
	if err := where.Find(&releases).Error; err != nil {
		return nil, err
	}
	return groupReleases(releases), nil
}

func (mgr *AppDatabaseMgr) FindReleases(condition map[string]interface{}) ([]*Release, error) {
	releases := []*Release{}
	err := mgr.db.Where(condition).Find(&releases).Error
	if err != nil {
		return nil, err
	}
	return releases, nil
}

func (mgr *AppDatabaseMgr) UpdateRelease(release *Release) error {
	return mgr.db.Save(release).Error
}
