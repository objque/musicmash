package db

import "time"

type InternalNotification struct {
	ArtistID   int64     `db:"artist_id"`
	ArtistName string    `db:"artist_name"`
	ReleaseID  uint64    `db:"release_id"`
	Title      string    `db:"title"`
	Released   time.Time `db:"released"`
	StoreID    string    `db:"store_id"`
	StoreName  string    `db:"store_name"`
	Poster     string    `db:"poster"`
	Type       string    `db:"type"`
	UserName   string    `db:"user_name"`
	Service    *string   `db:"service"`
	Data       *string   `db:"data"`
	Explicit   bool      `db:"explicit"`
}

func (r *InternalNotification) IsReleaseComing() bool {
	// if release day tomorrow or later, than that means coming release is here
	return r.Released.After(time.Now().UTC().Truncate(24 * time.Hour))
}

func (mgr *AppDatabaseMgr) FindNotReceivedNotifications() ([]*InternalNotification, error) {
	const query = `
SELECT releases.id as release_id,
       releases.artist_id,
       artists.name as artist_name,
       releases.title,
       releases.released,
       releases.store_id,
       releases.store_name,
       releases.poster,
       releases.type,
       releases.explicit,
       subscriptions.user_name,
       notification_settings.service,
       notification_settings.data
FROM releases
INNER JOIN subscriptions ON subscriptions.artist_id=releases.artist_id
              AND subscriptions.created_at <= releases.released
LEFT JOIN artists ON releases.artist_id=artists.id
LEFT JOIN notification_settings ON notification_settings.user_name=subscriptions.user_name
WHERE NOT EXISTS (
       SELECT user_name, release_id, is_coming
       FROM notifications
       WHERE notifications.user_name=subscriptions.user_name
              AND notifications.release_id=releases.id
              AND notifications.is_coming=(datetime(releases.released) > datetime($1))
)
ORDER BY subscriptions.user_name;
`

	future := time.Now().UTC().Truncate(24 * time.Hour)
	notifications := []*InternalNotification{}
	err := mgr.newdb.Select(&notifications, query, future.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	return notifications, nil
}
