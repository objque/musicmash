package notifier

import (
	"strconv"
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/musicmash/musicmash/internal/notifier/telegram"
	"github.com/pkg/errors"
)

func Notify() {
	last, err := db.DbMgr.GetLastActionDate(db.ActionNotify)
	if err != nil {
		log.Error(errors.Wrap(err, "tried to get last_action for notify stage"))
		return
	}

	NotifyWithPeriod(last.Date)
}

func markReleaseAsDeliveredTo(userName string, releaseID uint64, isComing bool) error {
	return db.DbMgr.CreateNotification(&db.Notification{
		Date:      time.Now().UTC(),
		UserName:  userName,
		ReleaseID: releaseID,
		IsComing:  isComing,
	})
}

func NotifyWithPeriod(period time.Time) {
	notifications, err := db.DbMgr.FindNotReceivedNotifications()
	if err != nil {
		log.Error(errors.Wrapf(err, "tried to notify users, but can't get new releases for date %v", period))
		return
	}

	if len(notifications) == 0 {
		log.Info("Not delivered notifications not found")
		return
	}

	for _, notification := range notifications {
		// todo: remove this after switching from gorm to sqlx
		notification.Release.ID = notification.ReleaseID
		notification.Release.Poster = notification.ReleasePoster
		message := makeMessage(notification.ArtistName, &notification.Release)
		message.ChatID, err = strconv.ParseInt(notification.Data, 10, 64)
		if err != nil {
			log.Warnf("user_name (%s) has broken %s data: '%v'",
				notification.UserName, notification.Service, notification.Data)
			continue
		}
		if err := telegram.SendMessage(message); err != nil {
			log.Error(errors.Wrapf(err, "tried to send message into telegram chat with id %v", notification.ReleaseID))
			continue
		}

		_ = markReleaseAsDeliveredTo(notification.UserName, notification.ReleaseID, notification.IsComing())
	}
}
