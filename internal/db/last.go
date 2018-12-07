package db

import "time"

const (
	ActionFetch   = "fetch"
	ActionReFetch = "refetch"
	ActionNotify  = "notify"
)

type LastAction struct {
	ID     int `gorm:"primary_key" sql:"AUTO_INCREMENT"`
	Date   time.Time
	Action string
}

type LastActionMgr interface {
	GetLastActionDate(action string) (*LastAction, error)
	SetLastActionDate(action string, time time.Time) error
}

func (mgr *AppDatabaseMgr) GetLastActionDate(action string) (*LastAction, error) {
	last := LastAction{}
	if err := mgr.db.First(&last, "action = ?", action).Error; err != nil {
		return nil, err
	}
	return &last, nil
}

func (mgr *AppDatabaseMgr) SetLastActionDate(action string, time time.Time) error {
	last, err := mgr.GetLastActionDate(action)
	if err != nil {
		return mgr.db.Create(&LastAction{
			Action: action,
			Date:   time,
		}).Error
	}

	last.Date = time
	return mgr.db.Save(last).Error
}
