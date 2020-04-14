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
	last, err := db.Mgr.GetLastActionDate(db.ActionNotify)
	if err != nil {
		log.Error(errors.Wrap(err, "tried to get last_action for notify stage"))
		return
	}

	NotifyWithPeriod(last.Date)
}

func markReleaseAsDelivered(userName string, releaseID uint64, isComing bool) error {
	return db.Mgr.CreateNotification(&db.Notification{
		Date:      time.Now().UTC(),
		UserName:  userName,
		ReleaseID: releaseID,
		IsComing:  isComing,
	})
}

func NotifyWithPeriod(period time.Time) {
	notifications, err := db.Mgr.FindNotReceivedNotifications()
	if err != nil {
		log.Error(errors.Wrapf(err, "tried to notify users, but can't get new releases for date %v", period))
		return
	}

	if len(notifications) == 0 {
		log.Info("Not delivered notifications not found")
		return
	}

	for _, notification := range notifications {
		message := makeMessage(notification.ArtistName, notification)
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

		_ = markReleaseAsDelivered(notification.UserName, notification.ReleaseID, notification.IsReleaseComing())
	}
}
