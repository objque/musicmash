package itunes

import (
	"strconv"
	"time"

	"github.com/objque/musicmash/internal/clients/itunes"
	"github.com/objque/musicmash/internal/clients/itunes/albums"
	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
	"github.com/pkg/errors"
)

const (
	posterWidth  = 500
	posterHeight = 500
)

func isLatest(album *albums.Album) bool {
	now := time.Now().UTC().Truncate(time.Hour * 24)
	yesterday := now.Add(-time.Hour * 48)
	return album.Attributes.ReleaseDate.Value.UTC().After(yesterday)
}

type Fetcher struct {
	Provider *itunes.Provider
}

func (f *Fetcher) GetStoreName() string {
	return "itunes"
}

func (f *Fetcher) fetchWorker(id int, artists <-chan *db.ArtistStoreInfo, done chan<- int) {
	for artist := range artists {
		artistID, err := strconv.ParseUint(artist.StoreID, 10, 64)
		if err != nil {
			log.Errorf("can't parse uint64 from '%s'", artist.StoreID)
			continue
		}

		release, err := albums.GetLatestArtistAlbum(f.Provider, artistID)
		if err != nil {
			if err == albums.AlbumsNotFoundErr {
				log.Debugf("Artist '%s' with id %s hasn't albums", artist.ArtistName, artist.StoreID)
				continue
			}

			log.Error(errors.Wrapf(err, "tried to get albums for '%s' with id %d", artist.ArtistName, artist.StoreID))
			continue
		}

		if !isLatest(release) {
			continue
		}

		err = db.DbMgr.EnsureReleaseExists(&db.Release{
			// TODO (m.kalinin): remove - Single, -Ep prefixes
			StoreName:  f.GetStoreName(),
			StoreID:    release.ID,
			ArtistName: artist.ArtistName,
			Title:      release.Attributes.Name,
			Poster:     release.Attributes.Artwork.GetLink(posterWidth, posterHeight),
			Released:   release.Attributes.ReleaseDate.Value,
		})
		if err != nil {
			log.Errorf("can't save release from '%s' with id '%s': %v", f.GetStoreName(), release.ID, err)
		}
	}
	done <- id
}

func (f *Fetcher) FetchAndSave(done chan<- bool) {
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

	done <- true
}
