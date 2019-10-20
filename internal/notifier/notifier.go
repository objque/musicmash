package notifier

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
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
	releases, err := db.DbMgr.FindNewReleases(period)
	if err != nil {
		log.Error(errors.Wrapf(err, "tried to notify users, but can't get new releases for date '%v'", period))
		return
	}

	artists := map[int64]*db.Artist{}
	for _, release := range releases {
		chats, err := db.DbMgr.GetAllChatsThatSubscribedFor(release.ArtistID)
		switch err {
		case gorm.ErrRecordNotFound:
			log.Debugf("No one subscribed for '%s'", release.ArtistID)
			continue
		case nil:
			break
		default:
			log.Error(err)
			continue
		}

		if len(chats) == 0 {
			continue
		}

		artist, ok := artists[release.ArtistID]
		if !ok {
			artist, _ = db.DbMgr.GetArtistWithFullInfo(release.ArtistID)
			artists[artist.ID] = artist
		}
		for _, chat := range chats {
			_, err := db.DbMgr.IsUserNotified(chat.UserName, release.ID, release.IsComing())
			switch err {
			case nil:
				log.Debugln(fmt.Sprintf("user '%s' already notified about '%d'", chat.UserName, release.ID))
				continue
			case gorm.ErrRecordNotFound:
				break
			default:
				log.Error(err)
				continue
			}

			if err := telegram.SendMessage(makeMessage(artist.Name, release)); err != nil {
				log.Error(errors.Wrapf(err, "tried to send message into telegram chat with id '%d'", chat.ID))
				continue
			}
			_ = markReleaseAsDeliveredTo(chat.UserName, release.ID, release.IsComing())
		}
	}
}
