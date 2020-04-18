package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Subscription struct {
	ID        uint64    `json:"-" db:"id"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UserName  string    `json:"-" db:"user_name"`
	ArtistID  int64     `json:"artist_id" db:"artist_id"`
}

type SubscriptionMgr interface {
	GetUserSubscriptions(userName string) ([]*Subscription, error)
	CreateSubscription(subscription *Subscription) error
	SubscribeUser(userName string, artists []int64) error
	UnSubscribeUser(userName string, artists []int64) error
}

func (mgr *AppDatabaseMgr) CreateSubscription(subscription *Subscription) error {
	const query = "insert into subscriptions (created_at, user_name, artist_id) VALUES (?, ?, ?)"

	now := subscription.CreatedAt.Format("2006-01-02T15:04:05")
	result, err := mgr.newdb.Exec(query, now, subscription.UserName, subscription.ArtistID)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	subscription.ID = uint64(id)
	return nil
}

func (mgr *AppDatabaseMgr) GetUserSubscriptions(userName string) ([]*Subscription, error) {
	const query = "select * from subscriptions where user_name = ?"

	subs := []*Subscription{}
	if err := mgr.newdb.Select(&subs, query, userName); err != nil {
		return nil, err
	}

	return subs, nil
}

func (mgr *AppDatabaseMgr) SubscribeUser(userName string, artists []int64) error {
	const rawquery = `
insert into subscriptions (created_at, user_name, artist_id)
select ? as created_at, ? as user_name, id as artist_id from artists
where
    artist_id in (?) and
    artist_id not in (select artist_id from subscriptions where user_name = ?)`

	now := time.Now().UTC().Format("2006-01-02T15:04:05")
	query, args, err := sqlx.In(rawquery, now, userName, artists, userName)
	if err != nil {
		return err
	}

	_, err = mgr.newdb.Exec(query, args...)
	return err
}

func (mgr *AppDatabaseMgr) UnSubscribeUser(userName string, artists []int64) error {
	const rawquery = "delete from subscriptions where user_name = ? and artist_id in (?)"

	query, args, err := sqlx.In(rawquery, userName, artists)
	if err != nil {
		return err
	}

	_, err = mgr.newdb.Exec(query, args...)
	return err
}
