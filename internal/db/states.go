package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type SubscriptionState string

const (
	ProcessingState SubscriptionState = "Processing"
	CompleteState   SubscriptionState = "Complete"
)

type State struct {
	ID        string            `gorm:"primary_key" json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Value     SubscriptionState `json:"value"`
}

type StateMgr interface {
	GetState(id string) (*State, error)
	UpdateState(id string, state SubscriptionState) error
}

func (mgr *AppDatabaseMgr) GetState(id string) (*State, error) {
	state := State{}
	if err := mgr.db.Where("id = ?", id).Find(&state).Error; err != nil {
		return nil, err
	}
	return &state, nil
}

func (mgr *AppDatabaseMgr) UpdateState(id string, state SubscriptionState) error {
	dbState, err := mgr.GetState(id)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}

		return mgr.db.Create(&State{ID: id, Value: state}).Error
	}
	return mgr.db.Model(dbState).Update("value", state).Error
}
