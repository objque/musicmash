package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Release struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"-"`
	ArtistID  int64     `json:"artist_id"`
	Title     string    `json:"title" gorm:"size:1000"`
	Poster    string    `json:"poster"`
	Released  time.Time `gorm:"type:datetime" json:"released"`
	StoreName string    `gorm:"unique_index:idx_rel_store_name_store_id" json:"store_name"`
	StoreID   string    `gorm:"unique_index:idx_rel_store_name_store_id" json:"store_id"`
}

type ReleaseMgr interface {
	EnsureReleaseExists(release *Release) error
	GetAllReleases() ([]*Release, error)
	FindReleases(condition map[string]interface{}) ([]*Release, error)
	FindNewReleases(date time.Time) ([]*Release, error)
	FindArtistsWithNewReleases(date time.Time) ([]int64, error)
	UpdateRelease(release *Release) error
}

func (r *Release) IsComing() bool {
	// if release day tomorrow or later, than that means coming release is here
	return r.Released.After(time.Now().UTC().Truncate(24 * time.Hour))
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
	var releases = []*Release{}
	return releases, mgr.db.Find(&releases).Error
}

func (mgr *AppDatabaseMgr) FindNewReleases(date time.Time) ([]*Release, error) {
	releases := []*Release{}
	if err := mgr.db.Where("created_at >= ?", date).Find(&releases).Error; err != nil {
		return nil, err
	}
	return releases, nil
}

func (mgr *AppDatabaseMgr) FindArtistsWithNewReleases(date time.Time) ([]int64, error) {
	ids := []int64{}
	const query = "select artist_id from releases where created_at >= ? group by artist_id"
	if err := mgr.db.Raw(query, date).Pluck("artist_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
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
