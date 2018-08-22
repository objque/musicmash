package notifier

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
	"github.com/objque/musicmash/internal/notify"
	"github.com/pkg/errors"
)

const delay = time.Minute * 15

func Run() {
	for {
		last, err := db.DbMgr.GetLastActionDate(db.ActionNotify)
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				log.Warnln("can't find last notify date. set as now()")
				db.DbMgr.SetLastActionDate(db.ActionNotify, time.Now().UTC())
				time.Sleep(delay)
				continue
			}

			log.Error(errors.Wrap(err, "tried to notify users, but can't get last notify action"))
			time.Sleep(delay)
			continue
		}

		releases, err := db.DbMgr.FindNewReleases(last.Date)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to notify users, but can't get new releases for date '%v'", last.Date))
			time.Sleep(delay)
			continue
		}

		for _, release := range releases {
			chats, err := db.DbMgr.GetAllChatsThatSubscribedFor(release.ArtistName)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					log.Debugf("No one subscribed for '%s'", release.ArtistName)
					continue
				}

				log.Error(err)
				continue
			}

			for _, chat := range chats {
				err := notify.Service.Send(map[string]interface{}{
					"chatID":  chat.ID,
					"release": release,
				})
				if err != nil {
					log.Error(errors.Wrapf(err, "tried to send message into telegram chat with id '%d'", chat.ID))
				}
			}
		}

		db.DbMgr.SetLastActionDate(db.ActionNotify, time.Now().UTC())
		time.Sleep(delay)
	}
}
