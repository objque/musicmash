package v2

import (
	"strconv"
	"time"

	itunes "github.com/objque/musicmash/internal/clients/itunes/v2"
	"github.com/objque/musicmash/internal/clients/itunes/v2/albums"
	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/fetcher/handlers"
	"github.com/objque/musicmash/internal/log"
	"github.com/pkg/errors"
)

func isLatest(album *albums.Album) bool {
	now := time.Now().UTC().Truncate(time.Hour * 24)
	yesterday := now.Add(-time.Hour * 48)
	return album.Attributes.ReleaseDate.Value.UTC().After(yesterday)
}

type Fetcher struct {
	handlers []handlers.StoreHandler
	Provider *itunes.Provider
}

func (f *Fetcher) fetchWorker(id int, artists <-chan *db.Artist, releases chan<- *db.Release, done chan<- int) {
	for artist := range artists {
		release, err := albums.GetLatestArtistAlbum(f.Provider, artist.StoreID)
		if err != nil {
			log.Error(err)
			continue
		}

		if !isLatest(release) {
			continue
		}

		releaseID, _ := strconv.ParseUint(release.ID, 10, 64)
		if db.DbMgr.IsReleaseExists(releaseID) {
			continue
		}

		dbRelease := db.Release{
			ArtistName: artist.Name,
			Date:       release.Attributes.ReleaseDate.Value,
			StoreID:    releaseID,
		}
		err = db.DbMgr.CreateRelease(&dbRelease)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to save release with id %v", release.ID))
			continue
		}

		releases <- &dbRelease
	}
	done <- id
}

func (f *Fetcher) fetch() ([]*db.Release, error) {
	// load all artists from the db
	artists, err := db.DbMgr.GetAllArtists()
	if err != nil {
		return nil, errors.Wrap(err, "can't load artists from the db")
	}

	jobs := make(chan *db.Artist, len(artists))
	releases := make(chan *db.Release, len(artists))
	done := make(chan int, config.Config.Fetching.Workers)

	// Starts up X workers, initially blocked because there are no jobs yet.
	for w := 1; w <= config.Config.Fetching.Workers; w++ {
		go f.fetchWorker(w, jobs, releases, done)
	}

	// Here we send `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for _, id := range artists {
		jobs <- id
	}
	close(jobs)

	for w := 1; w <= config.Config.Fetching.Workers; w++ {
		log.Debugf("#%d fetch-worker done", <-done)
	}
	close(releases)
	close(done)

	result := []*db.Release{}
	for release := range releases {
		result = append(result, release)
	}
	return result, nil
}

func (f *Fetcher) FetchAndProcess() error {
	releases, err := f.fetch()
	if err != nil {
		return err
	}

	for _, handler := range f.handlers {
		go handler.Fetch(releases)
	}
	return nil
}

func (f *Fetcher) RegisterHandler(handler handlers.StoreHandler) {
	f.handlers = append(f.handlers, handler)
}
