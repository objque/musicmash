package db

import (
	"time"
)

type LastFetch struct {
	ID   int32 `gorm:"primary_key"`
	Date time.Time
}

type LastFetchMgr interface {
	GetLastFetch() (*LastFetch, error)
	SetLastFetch(time time.Time) error
}

func (mgr *AppDatabaseMgr) GetLastFetch() (*LastFetch, error) {
	last := LastFetch{}
	if err := mgr.db.Last(&last).Error; err != nil {
		return nil, err
	}
	return &last, nil
}

func (mgr *AppDatabaseMgr) SetLastFetch(time time.Time) error {
	last, err := mgr.GetLastFetch()
	if err != nil {
		return mgr.db.Create(&LastFetch{
			Date: time,
		}).Error
	}

	last.Date = time
	return mgr.db.Save(last).Error
}
