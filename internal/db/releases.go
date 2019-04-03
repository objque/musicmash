package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/log"
)

type ReleaseStore struct {
	StoreURL  string `json:"url"`
	StoreName string `json:"name"`
	StoreID   string `json:"id"`
}

func (r *ReleaseStore) MarshalJSON() ([]byte, error) {
	s := struct {
		StoreURL  string `json:"url"`
		StoreName string `json:"name"`
		StoreID   string `json:"id"`
	}{
		StoreName: r.StoreName,
		StoreID:   r.StoreID,
	}

	if store, ok := config.Config.Stores[r.StoreName]; ok {
		s.StoreURL = fmt.Sprintf(store.ReleaseURL, r.StoreID)
	} else {
		log.Warnf("release_url for '%s' missed in config. User will receive empty link", r.StoreName)
	}
	return json.Marshal(&s)
}

type Release struct {
	ID         uint64          `json:"-"`
	CreatedAt  time.Time       `json:"-"`
	ArtistName string          `json:"artist_name"`
	Title      string          `json:"title" gorm:"size:1000"`
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
	FindArtistRecentReleases(artistName string) ([]*Release, error)
	FindArtistAnnouncedReleases(artistName string) ([]*Release, error)
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
	var releases = []*Release{}
	return releases, mgr.db.Find(&releases).Error
}

func (mgr *AppDatabaseMgr) GetReleasesForUserFilterByPeriod(userName string, since, till time.Time) ([]*Release, error) {
	// inner query: select artist_name from subscriptions where user_name = XXX
	// select * from releases where artist_name in (INNER) and released >= ? and released <= ? order by released desc
	releases := []*Release{}
	const query = "select artist_name from subscriptions where user_name = ?"
	innerQuery := mgr.db.Raw(query, userName).QueryExpr()
	where := mgr.db.Where("artist_name in (?) and released >= ? and released <= ?", innerQuery, since, till)
	// order by desc means first will be today releases
	if err := where.Order("title").Find(&releases).Error; err != nil {
		return nil, err
	}
	return releases, nil
}

func (mgr *AppDatabaseMgr) GetReleasesForUserSince(userName string, since time.Time) ([]*Release, error) {
	// inner query: select artist_name from subscriptions where user_name = XXX
	// select * from releases where artist_name in (INNER) and released >= ? order by released desc
	releases := []*Release{}
	const query = "select artist_name from subscriptions where user_name = ?"
	innerQuery := mgr.db.Raw(query, userName).QueryExpr()
	where := mgr.db.Where("artist_name in (?) and released >= ?", innerQuery, since)
	// order by asc means first will be releases that will be available very soon
	if err := where.Order("title").Find(&releases).Error; err != nil {
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
	releases := []*Release{}
	const query = "" +
		// releases from artists that user follow
		"artist_name IN (SELECT artist_name FROM subscriptions WHERE user_name = ?) " +
		// announced and fresh releases
		"AND (created_at >= ? OR released = ?) " +
		// not delivered releases
		"AND releases.id NOT IN (" +
		"	SELECT release_id FROM notifications WHERE user_name = ? and date like ?" +
		")"
	// gorm inserts str arguments with additional quotes, so we must provide date with percentages
	// to avoid this part of query: "and date like '%'2019-04-03'%'"
	likeCondition := fmt.Sprintf("%%%s%%", date.Format("2006-01-02"))
	where := mgr.db.Where(query, userName, date, date.Truncate(time.Hour*24), userName, likeCondition)
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

func (mgr *AppDatabaseMgr) FindArtistRecentReleases(artistName string) ([]*Release, error) {
	releases := []*Release{}
	err := mgr.db.Where("artist_name = ? and released <= ?", artistName, time.Now().UTC()).Find(&releases).Error
	if err != nil {
		return nil, err
	}
	return releases, nil
}

func (mgr *AppDatabaseMgr) FindArtistAnnouncedReleases(artistName string) ([]*Release, error) {
	releases := []*Release{}
	err := mgr.db.Where("artist_name = ? and released > ?", artistName, time.Now().UTC()).Find(&releases).Error
	if err != nil {
		return nil, err
	}
	return releases, nil
}
