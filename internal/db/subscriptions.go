package db

import (
	"github.com/objque/musicmash/internal/log"
)

type Subscription struct {
	ID         int64  `gorm:"primary_key"`
	UserID     string `sql:"index" gorm:"unique_index:idx_user_id_artist_name"`
	ArtistName string `gorm:"unique_index:idx_user_id_artist_name"`
}

type SubscriptionMgr interface {
	IsUserSubscribedForArtist(userID, artistName string) bool
	FindAllUserSubscriptions(userID string) ([]*Subscription, error)
	EnsureSubscriptionExists(subscription *Subscription) error
	SubscribeUserForArtists(userID string, artists []string) error
	UnsubscribeUserFromArtists(userID string, artist []string) error
}

func (mgr *AppDatabaseMgr) IsUserSubscribedForArtist(userID, artistName string) bool {
	subscription := Subscription{}
	if err := mgr.db.Where("user_id = ? and artist_name = ?", userID, artistName).First(&subscription).Error; err != nil {
		return false
	}

	return true
}

func (mgr *AppDatabaseMgr) FindAllUserSubscriptions(userID string) ([]*Subscription, error) {
	subscriptions := []*Subscription{}
	if err := mgr.db.Where("user_id = ?", userID).Find(&subscriptions).Error; err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (mgr *AppDatabaseMgr) EnsureSubscriptionExists(subscription *Subscription) error {
	if !mgr.IsUserSubscribedForArtist(subscription.UserID, subscription.ArtistName) {
		return mgr.db.Create(subscription).Error
	}
	return nil
}

func (mgr *AppDatabaseMgr) SubscribeUserForArtists(userID string, artists []string) error {
	const sql = `insert into subscriptions (user_id, artist_name) values (?, ?)`
	for i := range artists {
		if err := mgr.db.Exec(sql, userID, artists[i]).Error; err != nil {
			log.Error(err)
		}
	}
	return nil
}

func (mgr *AppDatabaseMgr) UnsubscribeUserFromArtists(userID string, artists []string) error {
	const sql = `delete from subscriptions where user_id = ? and artist_name in (?)`
	return mgr.db.Exec(sql, userID, artists).Error
}
