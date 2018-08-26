package yandex

import (
	"strconv"

	"github.com/objque/musicmash/internal/clients/itunes"
	"github.com/objque/musicmash/internal/clients/yandex"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
	"github.com/pkg/errors"
)

type YandexHandler struct {
	api *yandex.Client
}

func New(url string) *YandexHandler {
	return &YandexHandler{
		api: yandex.New(url),
	}
}

func (h *YandexHandler) Fetch(releases []*db.Release) {
	for _, dbRelease := range releases {
		release, err := itunes.Lookup(dbRelease.StoreID)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to lookup '%d' before searching release in yandex", dbRelease.StoreID))
			continue
		}

		yandexID, err := find(h.api, release.ArtistName, release.CollectionName)
		if err != nil {
			if err == ArtistNotFoundErr || err == ReleaseNotFoundErr {
				continue
			}

			log.Error(err)
			continue
		}

		log.Infof("Found a new info from '%s' in yandex.music: '%d'", dbRelease.ArtistName, yandexID)
		err = db.DbMgr.EnsureReleaseExistsInStore(h.GetStoreName(), strconv.Itoa(yandexID), dbRelease.ID)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to save release in %s with id '%v'", h.GetStoreName(), dbRelease.ID))
			continue
		}
	}
}

func (h *YandexHandler) GetStoreName() string {
	return "yandex"
}
