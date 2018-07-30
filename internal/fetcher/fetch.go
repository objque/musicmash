package fetcher

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/itunes"
	"github.com/objque/musicmash/internal/log"
	"github.com/pkg/errors"
)

func saveIfNewestRelease(artist string, release *itunes.LastRelease) bool {
	if !release.IsLatest() {
		return false
	}

	log.Infof("Found a new release from '%s': '%d'", artist, release.ID)
	db.DbMgr.EnsureReleaseExists(&db.Release{
		// NOTE (m.kalinin): we provide artist because if artist releases feat with someone, then
		// release will contain incorrect name.
		ArtistName: artist,
		StoreURL:   release.URL,
	})
	return true
}

func fetch() error {
	// load all artists from the db
	artists, err := db.DbMgr.GetAllArtists()
	if err != nil {
		return errors.Wrap(err, "can't load artists from the db")
	}

	// load releases from the store
	for _, artist := range artists {
		releaseInfo, err := itunes.GetArtistInfo(artist.StoreID)
		if err != nil {
			log.Error(err)
			continue
		}

		saveIfNewestRelease(artist.Name, releaseInfo)
	}
	return nil
}

func isMustFetch() bool {
	last, err := db.DbMgr.GetLastFetch()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		}

		log.Error(err)
		return false
	}

	log.Debugf("Last fetch was at '%s'", last.Date.String())
	return calcDiffHours(last.Date) > config.Config.Fetching.CountOfSkippedHoursToFetch
}

func Run() {
	for {
		if isMustFetch() {
			now := time.Now().UTC()
			log.Infof("Start fetching stage for '%s'...", now.String())
			if err := fetch(); err != nil {
				log.Error(err)
			} else {
				log.Infof("Finish fetching stage '%s'...", time.Now().UTC().String())
				db.DbMgr.SetLastFetch(now)
			}
		}

		time.Sleep(time.Hour)
	}
}
