package itunes

import (
	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/fetcher/handlers"
	"github.com/objque/musicmash/internal/itunes"
	"github.com/objque/musicmash/internal/log"
	"github.com/objque/musicmash/internal/notify"
	"github.com/pkg/errors"
)

const storeName = "itunes"

type AppleMusicHandler struct{}

func (h *AppleMusicHandler) GetStoreName() string {
	return storeName
}

func (h *AppleMusicHandler) Fetch(releases []*itunes.LastRelease) []*handlers.ReleaseData {
	data := []*handlers.ReleaseData{}
	for _, release := range releases {
		if !release.IsComing {
			log.Infof("Found a new release from '%s': '%d'", release.ArtistName, release.ID)
		} else {
			log.Infof("Found a new pre-release from '%s': '%d'", release.ArtistName, release.ID)
		}

		err := db.DbMgr.CreateRelease(&db.Release{
			ArtistName: release.ArtistName,
			Date:       release.Date,
			StoreID:    release.ID,
			StoreType:  storeName,
		})
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to save release for %s %v", storeName, release.ID))
			continue
		}
		data = append(data, &handlers.ReleaseData{
			ArtistName: release.ArtistName,
			Date:       release.Date,
			StoreID:    release.ID,
		})
	}
	return data
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
