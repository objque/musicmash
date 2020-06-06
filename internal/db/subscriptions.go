package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Subscription struct {
	ID           uint64    `json:"-"`
	CreatedAt    time.Time `json:"-"`
	UserName     string    `json:"-"             db:"user_name"`
	ArtistID     int64     `json:"artist_id"     db:"artist_id"`
	ArtistName   string    `json:"artist_name"   db:"artist_name"`
	ArtistPoster string    `json:"artist_poster" db:"artist_poster"`
}

func (mgr *AppDatabaseMgr) CreateSubscription(subscription *Subscription) error {
	const query = "insert into subscriptions (created_at, user_name, artist_id) VALUES ($1, $2, $3)"

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
	const query = "" +
		"select" +
		" user_name," +
		" artists.id as artist_id," +
		" artists.name as artist_name," +
		" artists.poster as artist_poster " +
		"from subscriptions " +
		"left join artists on subscriptions.artist_id=artists.id " +
		"where subscriptions.user_name = $1"

	subs := []*Subscription{}
	if err := mgr.newdb.Select(&subs, query, userName); err != nil {
		return nil, err
	}

	return subs, nil
}

func (mgr *AppDatabaseMgr) SubscribeUser(userName string, artists []int64) error {
	const rawquery = `
insert into subscriptions (created_at, user_name, artist_id)
select $1 as created_at, $2 as user_name, id as artist_id from artists
where
    artist_id in ($3) and
    artist_id not in (select artist_id from subscriptions where user_name = $4)`

	now := time.Now().UTC().Format("2006-01-02T15:04:05")
	query, args, err := sqlx.In(rawquery, now, userName, artists, userName)
	if err != nil {
		return err
	}

	_, err = mgr.newdb.Exec(query, args...)
	return err
}

func (mgr *AppDatabaseMgr) UnSubscribeUser(userName string, artists []int64) error {
	const rawquery = "delete from subscriptions where user_name = $1 and artist_id in ($2)"

	query, args, err := sqlx.In(rawquery, userName, artists)
	if err != nil {
		return err
	}

	_, err = mgr.newdb.Exec(query, args...)
	return err
}
