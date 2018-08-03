package fetcher

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/itunes"
	"github.com/objque/musicmash/internal/log"
	"github.com/objque/musicmash/internal/notify"
	"github.com/pkg/errors"
)

func saveIfNewestRelease(artist string, release *itunes.LastRelease) bool {
	if !release.IsLatest() {
		return false
	}

	if db.DbMgr.IsReleaseExists(release.ID) {
		return false
	}

	if release.IsComing {
		log.Infof("Found pre-release from '%s'", artist)
		return false
	}

	log.Infof("Found a new release from '%s': '%d'", artist, release.ID)
	db.DbMgr.CreateRelease(&db.Release{
		ArtistName: artist,
		Date:       release.Date,
		StoreID:    release.ID,
	})
	notify.Service.Send(map[string]interface{}{
		"chatID":    int64(35152258),
		"releaseID": release.ID,
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
			if err == itunes.ArtistInactiveErr {
				log.Debugln(errors.Wrapf(err, "artist: '%s'#%d", artist.Name, artist.StoreID))
				continue
			}

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
	diff := calcDiffHours(last.Date)
	log.Debugf("Diff between hours: %v", diff)
	return diff > config.Config.Fetching.CountOfSkippedHoursToFetch
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
			log.Debugf("Elapsed time '%s'", time.Now().UTC().Sub(now).String())
		}

		time.Sleep(time.Hour)
	}
}
