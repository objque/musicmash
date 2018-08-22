package itunes

import (
	"strconv"

	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
	"github.com/pkg/errors"
)

type AppleMusicHandler struct{}

func (h *AppleMusicHandler) GetStoreName() string {
	return "itunes"
}

func (h *AppleMusicHandler) Fetch(releases []*db.Release) {
	for _, release := range releases {
		log.Infof("Found a new info from '%s': '%d'", release.ArtistName, release.ID)
		err := db.DbMgr.EnsureReleaseExistsInStore(h.GetStoreName(), strconv.FormatUint(release.StoreID, 10), release.ID)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to save release in %s with id '%v'", h.GetStoreName(), release.ID))
			continue
		}
	}
}
