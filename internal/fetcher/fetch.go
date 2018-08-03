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

func saveIfNewestRelease(release *itunes.LastRelease) bool {
	if !release.IsLatest() {
		return false
	}

	if db.DbMgr.IsReleaseExists(release.ID) {
		return false
	}

	if !release.IsComing {
		log.Infof("Found a new release from '%s': '%d'", release.ArtistName, release.ID)
	} else {
		log.Infof("Found a new pre-release from '%s': '%d'", release.ArtistName, release.ID)
	}

	db.DbMgr.CreateRelease(&db.Release{
		ArtistName: release.ArtistName,
		Date:       release.Date,
		StoreID:    release.ID,
	})

	notify.Service.Send(map[string]interface{}{
		"chatID":          int64(35152258),
		"releaseID":       release.ID,
		"isFutureRelease": release.IsComing,
	})
	return true
}

func fetchWorker(id int, artists <-chan *db.Artist, releases chan<- *itunes.LastRelease, done chan<- int) {
	for artist := range artists {
		release, err := itunes.GetArtistInfo(artist.StoreID)
		if err != nil {
			if err == itunes.ArtistInactiveErr {
				log.Debugln(errors.Wrapf(err, "artist: '%s'#%d", artist.Name, artist.StoreID))
				releases <- nil
				continue
			}

			log.Error(err)
			releases <- nil
			continue
		}

		release.ArtistName = artist.Name
		releases <- release
	}
	done <- id
}

func saveWorker(id int, releases <-chan *itunes.LastRelease, done chan<- int) {
	for release := range releases {
		if release == nil {
			continue
		}

		saveIfNewestRelease(release)
	}
	done <- id
}

func fetch() error {
	// load all artists from the db
	artists, err := db.DbMgr.GetAllArtists()
	if err != nil {
		return errors.Wrap(err, "can't load artists from the db")
	}

	jobs := make(chan *db.Artist, len(artists))
	releases := make(chan *itunes.LastRelease, len(artists))
	fetchWorkersDone := make(chan int, config.Config.Fetching.Workers)
	saveWorkersDone := make(chan int, config.Config.Fetching.Workers)

	// Starts up X workers, initially blocked because there are no jobs yet.
	for w := 1; w <= config.Config.Fetching.Workers; w++ {
		go fetchWorker(w, jobs, releases, fetchWorkersDone)
		go saveWorker(w, releases, saveWorkersDone)
	}

	// Here we send `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for _, id := range artists {
		jobs <- id
	}
	close(jobs)

	for w := 1; w <= config.Config.Fetching.Workers; w++ {
		log.Debugf("#%d fetch-worker done\n", <-fetchWorkersDone)
	}
	close(releases)

	for w := 1; w <= config.Config.Fetching.Workers; w++ {
		log.Debugf("#%d save-worker done\n", <-saveWorkersDone)
	}
	close(fetchWorkersDone)
	close(saveWorkersDone)
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
