package deezer

import (
	"strconv"
	"sync"
	"time"

	"github.com/musicmash/musicmash/internal/clients/deezer"
	"github.com/musicmash/musicmash/internal/clients/deezer/albums"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/pkg/errors"
)

func isLatest(album *albums.Album) bool {
	now := time.Now().UTC().Truncate(time.Hour * 24)
	yesterday := now.Add(-time.Hour * 48)
	return album.Released.Value.UTC().After(yesterday)
}

type Fetcher struct {
	Provider     *deezer.Provider
	FetchWorkers int
}

func NewService(url string, fetchWorkers int) *Fetcher {
	return &Fetcher{
		Provider:     deezer.NewProvider(url),
		FetchWorkers: fetchWorkers,
	}
}

func (f *Fetcher) GetStoreName() string {
	return "deezer"
}

func (f *Fetcher) fetchWorker(id int, artists <-chan *db.ArtistStoreInfo, done chan<- int) {
	for artist := range artists {
		artistID, err := strconv.Atoi(artist.StoreID)
		if err != nil {
			log.Errorf("can't parse int from '%s'", artist.StoreID)
			continue
		}

		release, err := albums.GetLatestArtistAlbum(f.Provider, artistID)
		if err != nil {
			if err == albums.ErrAlbumsNotFound {
				log.Debugf("Artist '%s' with id %s hasn't albums", artist.ArtistName, artist.StoreID)
				continue
			}

			log.Error(errors.Wrapf(err, "tried to get albums for '%s' with id %s", artist.ArtistName, artist.StoreID))
			continue
		}

		if !isLatest(release) {
			continue
		}

		err = db.DbMgr.EnsureReleaseExists(&db.Release{
			StoreName:  f.GetStoreName(),
			StoreID:    strconv.Itoa(release.ID),
			ArtistName: artist.ArtistName,
			Title:      release.Title,
			Poster:     release.Poster,
			Released:   release.Released.Value,
		})
		if err != nil {
			log.Errorf("can't save release from '%s' with id '%s': %v", f.GetStoreName(), release.ID, err)
		}
	}
	done <- id
}

func (f *Fetcher) FetchAndSave(wg *sync.WaitGroup) {
	defer wg.Done()
	// load all artists from the db
	artists, err := db.DbMgr.GetArtistsForStore(f.GetStoreName())
	if err != nil {
		log.Error(errors.Wrap(err, "can't load artists from the db"))
		return
	}

	jobs := make(chan *db.ArtistStoreInfo, len(artists))
	_done := make(chan int, f.FetchWorkers)

	// Starts up X workers, initially blocked because there are no jobs yet.
	for w := 1; w <= f.FetchWorkers; w++ {
		go f.fetchWorker(w, jobs, _done)
	}

	// Here we send `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for _, id := range artists {
		jobs <- id
	}
	close(jobs)

	for w := 1; w <= f.FetchWorkers; w++ {
		log.Debugf("#%d fetch-worker wg", <-_done)
	}
	close(_done)
}

func (f *Fetcher) ReFetchAndSave(wg *sync.WaitGroup) {
	defer wg.Done()
	releases, err := db.DbMgr.FindReleases(map[string]interface{}{
		"poster":     "",
		"store_name": f.GetStoreName(),
	})
	if err != nil {
		log.Error(errors.Wrap(err, "tried to get releases for refetch from deezer"))
		return
	}

	for _, release := range releases {
		log.Debugf("trying to refetch %s release with id %s", f.GetStoreName(), release.StoreID)
		albumID, err := strconv.Atoi(release.StoreID)
		if err != nil {
			log.Errorf("can't cast str to int: '%s'", release.StoreID)
			continue
		}

		album, err := albums.GetByID(f.Provider, albumID)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried get info about album %s from %s", release.StoreID, f.GetStoreName()))
			continue
		}

		if album.Poster != "" {
			log.Debugf("found missed album poster for release %s from %s", release.StoreID, f.GetStoreName())
			release.Poster = album.Poster
			if err := db.DbMgr.UpdateRelease(release); err != nil {
				log.Error(errors.Wrapf(err, "tried to update release without poster with id %s", release.StoreID))
			} else {
				log.Infof("poster for album with id %s from %s was updated", release.StoreID, f.GetStoreName())
			}
			continue
		}

		log.Debugf("release with id %s from %s still exists without poster", release.StoreID, f.GetStoreName())
	}
}
