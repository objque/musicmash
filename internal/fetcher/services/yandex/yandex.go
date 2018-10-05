package yandex

import (
	"strconv"
	"sync"
	"time"

	"github.com/objque/musicmash/internal/clients/yandex"
	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
	"github.com/pkg/errors"
)

const (
	posterWidth  = 300
	posterHeight = 300
)

func isLatest(album *yandex.ArtistAlbum) bool {
	now := time.Now().UTC().Truncate(time.Hour * 24)
	yesterday := now.Add(-time.Hour * 48)
	return album.Released.Value.UTC().After(yesterday)
}

type Fetcher struct {
	API *yandex.Client
}

func NewService(url string) *Fetcher {
	return &Fetcher{
		API: yandex.New(url),
	}
}

func (f *Fetcher) GetStoreName() string {
	return "yandex"
}

func (f *Fetcher) fetchWorker(id int, artists <-chan *db.ArtistStoreInfo, done chan<- int) {
	for artist := range artists {
		artistID, err := strconv.Atoi(artist.StoreID)
		if err != nil {
			log.Errorf("can't parse int from '%s'", artist.StoreID)
			continue
		}

		album, err := f.API.GetArtistLatestAlbum(artistID)
		if err != nil {
			if err == yandex.AlbumsNotFoundErr {
				log.Debugf("Artist '%s' with id %d hasn't albums", artist.ArtistName, artist.StoreID)
				continue
			}

			log.Error(errors.Wrapf(err, "tried to get latest album for '%s' with id %s", artist.ArtistName, artist.StoreID))
			continue
		}

		if !isLatest(album) {
			continue
		}

		err = db.DbMgr.EnsureReleaseExists(&db.Release{
			StoreName:  f.GetStoreName(),
			StoreID:    strconv.Itoa(album.ID),
			ArtistName: artist.ArtistName,
			Title:      album.Title,
			Poster:     album.GetPosterWithSize(posterWidth, posterHeight),
			Released:   album.Released.Value,
		})
		if err != nil {
			log.Errorf("can't save release from '%s' with id '%s': %v", f.GetStoreName(), strconv.Itoa(album.ID), err)
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
	_done := make(chan int, config.Config.Fetching.Workers)

	// Starts up X workers, initially blocked because there are no jobs yet.
	// TODO (m.kalinin): replace with store-workers count
	for w := 1; w <= config.Config.Fetching.Workers; w++ {
		go f.fetchWorker(w, jobs, _done)
	}

	// Here we send `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for _, id := range artists {
		jobs <- id
	}
	close(jobs)

	for w := 1; w <= config.Config.Fetching.Workers; w++ {
		log.Debugf("#%d fetch-worker done", <-_done)
	}
	close(_done)
}
