package db

import "github.com/musicmash/musicmash/internal/log"

type Subscription struct {
	ID         int64  `json:"-" gorm:"primary_key"`
	UserName   string `json:"-" sql:"index" gorm:"unique_index:idx_user_name_artist_name"`
	ArtistName string `json:"artist_name" gorm:"unique_index:idx_user_name_artist_name"`
}

type SubscriptionMgr interface {
	IsUserSubscribedForArtist(userName, artistName string) bool
	FindAllUserSubscriptions(userName string) ([]*Subscription, error)
	EnsureSubscriptionExists(userName, artistName string) error
	SubscribeUserForArtists(userName string, artists []string) error
	UnsubscribeUserFromArtists(userName string, artist []string) error
}

func (mgr *AppDatabaseMgr) IsUserSubscribedForArtist(userName, artistName string) bool {
	subscription := Subscription{}
	if err := mgr.db.Where("user_name = ? and artist_name = ?", userName, artistName).First(&subscription).Error; err != nil {
		return false
	}

	return true
}

func (mgr *AppDatabaseMgr) FindAllUserSubscriptions(userName string) ([]*Subscription, error) {
	subscriptions := []*Subscription{}
	if err := mgr.db.Where("user_name = ?", userName).Find(&subscriptions).Error; err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (mgr *AppDatabaseMgr) EnsureSubscriptionExists(userName, artistName string) error {
	if !mgr.IsUserSubscribedForArtist(userName, artistName) {
		return mgr.db.Create(&Subscription{UserName: userName, ArtistName: artistName}).Error
	}
	return nil
}

func (mgr *AppDatabaseMgr) SubscribeUserForArtists(userName string, artists []string) error {
	const sql = `insert into subscriptions (user_name, artist_name) values (?, ?)`
	for i := range artists {
		if err := mgr.db.Exec(sql, userName, artists[i]).Error; err != nil {
			log.Error(err)
		}
	}
	return nil
}

func (mgr *AppDatabaseMgr) UnsubscribeUserFromArtists(userName string, artists []string) error {
	const sql = `delete from subscriptions where user_name = ? and artist_name in (?)`
	return mgr.db.Exec(sql, userName, artists).Error
}
