package db

import "time"

type InternalNotification struct {
	ArtistID   int64
	ArtistName string
	ReleaseID  uint64
	Title      string
	Released   time.Time
	StoreID    string
	StoreName  string
	Poster     string
	Type       string
	UserName   string
	Service    string
	Data       string
	Explicit   bool
}

type InternalNotificationMgr interface {
	FindNotReceivedNotifications() ([]*InternalNotification, error)
}

func (r *InternalNotification) IsReleaseComing() bool {
	// if release day tomorrow or later, than that means coming release is here
	return r.Released.After(time.Now().UTC().Truncate(24 * time.Hour))
}

func (mgr *AppDatabaseMgr) FindNotReceivedNotifications() ([]*InternalNotification, error) {
	future := time.Now().UTC().Truncate(24 * time.Hour)
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
              AND notifications.is_coming=(datetime(releases.released) > datetime(?))
)
ORDER BY subscriptions.user_name;
`

	notifications := []*InternalNotification{}
	err := mgr.db.Raw(query, future).Scan(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
