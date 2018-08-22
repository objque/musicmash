package itunes

import (
	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/fetcher/handlers"
	"github.com/objque/musicmash/internal/log"
	"github.com/objque/musicmash/internal/notify"
	"github.com/pkg/errors"
)

const storeName = "itunes"

type AppleMusicHandler struct{}

func (h *AppleMusicHandler) GetStoreName() string {
	return storeName
}

func (h *AppleMusicHandler) Fetch(releases []*db.Release) {
	for _, release := range releases {
		log.Infof("Found a new info from '%s': '%d'", release.ArtistName, release.ID)
		err := db.DbMgr.EnsureReleaseExistsInStore(storeName, string(release.ID), release.ID)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to save release in %s with id '%v'", storeName, release.ID))
			continue
		}
	}
}

func (h *AppleMusicHandler) NotifySubscribers(releases []*handlers.ReleaseData) {
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
				"chatID":    chat.ID,
				"releaseID": release.StoreID,
				"store":     storeName,
			})
			if err != nil {
				log.Error(errors.Wrapf(err, "tried to send message into telegram chat with id '%d'", chat.ID))
			}
		}
	}
}
