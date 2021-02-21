package db

import (
	"time"

	sq "github.com/Masterminds/squirrel"
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

type GetSubscriptionsOpts struct {
	Limit  *uint64
	Before *uint64
}

func applySubscriptionsFilters(query sq.SelectBuilder, opts *GetSubscriptionsOpts) sq.SelectBuilder {
	// OrderByClause method generates incorrect query and we can't pass ASC/DESC as an arg
	query = query.OrderBy("subscriptions.id DESC")

	if opts.Before != nil {
		query = query.Where("subscriptions.id < ?", *opts.Before)
	}

	if opts.Limit != nil {
		query = query.Limit(*opts.Limit)
	}

	return query
}

func (mgr *AppDatabaseMgr) CreateSubscription(subscription *Subscription) error {
	const query = "insert into subscriptions (created_at, user_name, artist_id) VALUES ($1, $2, $3) returning id"

	now := subscription.CreatedAt.Format("2006-01-02T15:04:05")

	return mgr.newdb.QueryRow(query, now, subscription.UserName, subscription.ArtistID).Scan(&subscription.ID)
}

func (mgr *AppDatabaseMgr) GetUserSubscriptions(userName string, opts *GetSubscriptionsOpts) ([]*Subscription, error) {
	query := sq.Select(
		"user_name",
		"artists.id as artist_id",
		"artists.name as artist_name",
		"artists.poster as artist_poster").
		From("subscriptions").
		LeftJoin("artists on subscriptions.artist_id=artists.id").
		Where("subscriptions.user_name = ?", userName)

	if opts != nil {
		query = applySubscriptionsFilters(query, opts)
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	subs := []*Subscription{}
	if err := mgr.newdb.Select(&subs, sql, args...); err != nil {
		return nil, err
	}

	return subs, nil
}

func (mgr *AppDatabaseMgr) SubscribeUser(userName string, artists []int64) error {
	const query = "INSERT INTO subscriptions (created_at, user_name, artist_id) " +
		"VALUES (now(), $1, $2) ON CONFLICT DO NOTHING"

	tx, err := mgr.Begin()
	if err != nil {
		return err
	}

	for _, artistID := range artists {
		if _, err = tx.newdb.Exec(query, userName, artistID); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (mgr *AppDatabaseMgr) UnSubscribeUser(userName string, artists []int64) error {
	const rawquery = "delete from subscriptions where user_name = ? and artist_id in (?)"

	query, args, err := sqlx.In(rawquery, userName, artists)
	if err != nil {
		return err
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	_, err = mgr.newdb.Exec(query, args...)
	return err
}
