package db

import "time"

type Subscription struct {
	ID        uint64    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UserName  string    `json:"-"`
	ArtistID  int64     `json:"artist_id"`
}

type SubscriptionMgr interface {
	GetUserSubscriptions(userName string) ([]*Subscription, error)
	GetArtistsSubscriptions(artists []int64) ([]*Subscription, error)
	CreateSubscription(subscription *Subscription) error
	SubscribeUser(userName string, artists []int64) error
	UnSubscribeUser(userName string, artists []int64) error
}

func (mgr *AppDatabaseMgr) CreateSubscription(subscription *Subscription) error {
	return mgr.db.Create(&subscription).Error
}

func (mgr *AppDatabaseMgr) GetUserSubscriptions(userName string) ([]*Subscription, error) {
	subs := []*Subscription{}
	err := mgr.db.Where("user_name = ?", userName).Find(&subs).Error
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (mgr *AppDatabaseMgr) GetArtistsSubscriptions(artists []int64) ([]*Subscription, error) {
	subs := []*Subscription{}
	err := mgr.db.Where("artist_id in (?)", artists).Order("user_name").Find(&subs).Error
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (mgr *AppDatabaseMgr) SubscribeUser(userName string, artists []int64) error {
	now := time.Now().UTC()
	const query = `
insert into subscriptions (created_at, user_name, artist_id)
select ? as created_at, ? as user_name, id as artist_id from artists
where
    artist_id in (?) and
    artist_id not in (select artist_id from subscriptions where user_name = ?)`

	return mgr.db.Exec(query, now, userName, artists, userName).Error
}

func (mgr *AppDatabaseMgr) UnSubscribeUser(userName string, artists []int64) error {
	const query = "delete from subscriptions where user_name = ? and artist_id in (?)"
	return mgr.db.Exec(query, userName, artists).Error
}
