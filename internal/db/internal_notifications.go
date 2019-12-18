package db

import "time"

type InternalNotification struct {
	Artist
	Release
	NotificationSettings
	ReleaseID     uint64
	ReleasePoster string
}

type InternalNotificationMgr interface {
	FindNotReceivedNotifications() ([]*InternalNotification, error)
}

func (mgr *AppDatabaseMgr) FindNotReceivedNotifications() ([]*InternalNotification, error) {
	future := time.Now().UTC().Truncate(24 * time.Hour)
	const query = `
SELECT releases.id as release_id,
       releases.artist_id,
       artists.name,
       releases.title,
       releases.released,
       releases.store_id,
       releases.store_name,
       releases.poster as release_poster,
       subscriptions.user_name,
       notification_settings.service,
       notification_settings.data
FROM releases
INNER JOIN subscriptions ON subscriptions.artist_id=releases.artist_id 
              AND subscriptions.created_at <= releases.released
LEFT JOIN notifications ON notifications.user_name=subscriptions.user_name
              AND notifications.release_id=releases.id
              AND is_coming=(releases.released >= ?)
LEFT JOIN artists ON releases.artist_id=artists.id
LEFT JOIN notification_settings ON notification_settings.user_name=subscriptions.user_name
WHERE notifications.user_name IS NULL
ORDER BY notifications.user_name
`

	notifications := []*InternalNotification{}
	err := mgr.db.Raw(query, future).Scan(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
