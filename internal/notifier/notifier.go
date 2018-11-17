package notifier

import (
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/musicmash/musicmash/internal/notifier/telegram"
	"github.com/pkg/errors"
)

func notify(chatID int64, releases []*db.Release) {
	for i := range releases {
		message := MakeMessage(releases[i])
		message.ChatID = chatID
		if err := telegram.SendMessage(message); err != nil {
			log.Error(errors.Wrapf(err, "tried to send release to '%d'", chatID))
		}
	}
}

func Notify() {
	last, err := db.DbMgr.GetLastActionDate(db.ActionNotify)
	if err != nil {
		log.Error(errors.Wrap(err, "tried to get last_action for notify stage"))
		return
	}

	users, err := db.DbMgr.GetUsersWithReleases(last.Date)
	if err != nil {
		log.Error(errors.Wrap(err, "tried to get users with releases for notify stage"))
		return
	}

	for _, user := range users {
		chat, err := db.DbMgr.FindChatByUserName(user)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to get chat for '%s' for notify stage", user))
			continue
		}

		// bug: method (FindNewReleasesForUser) receives albums that were founded after not truncated date
		// and misses releases that were announced sometime ago.
		//
		// also, at now we don't known which albums we sent user earlier.
		// so we may use workaround: notify only once per day (set count_of_skipped_hours=24 in the cfg for notifier)
		// and provide truncated date in the FindNewReleasesForUser method.
		releases, err := db.DbMgr.FindNewReleasesForUser(user, last.Date.Truncate(time.Hour*24))
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to get feed for '%s' for notify stage", user))
			return
		}

		notify(*chat, releases)
	}
}
