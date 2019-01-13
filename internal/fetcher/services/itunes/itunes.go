package itunes

import (
	"strconv"
	"strings"
	"sync"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/musicmash/musicmash/internal/clients/itunes/albums"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/pkg/errors"
)

const (
	posterWidth  = 500
	posterHeight = 500

	AlbumReleaseType  = " - Album"
	SingleReleaseType = " - Single"
	EPReleaseType     = " - EP"
	LPReleaseType     = " - LP"
)

type Fetcher struct {
	Provider     *itunes.Provider
	FetchWorkers int
}

func NewService(url string, fetchWorkers int, token string) *Fetcher {
	return &Fetcher{
		Provider:     itunes.NewProvider(url, token),
		FetchWorkers: fetchWorkers,
	}
}

func (f *Fetcher) GetStoreName() string {
	return "itunes"
}

func removeAlbumType(title string) string {
	title = strings.Replace(title, AlbumReleaseType, "", -1)
	title = strings.Replace(title, SingleReleaseType, "", -1)
	title = strings.Replace(title, EPReleaseType, "", -1)
	return strings.Replace(title, LPReleaseType, "", -1)
}

func (f *Fetcher) fetchWorker(id int, artists <-chan *db.ArtistStoreInfo, done chan<- int) {
	for artist := range artists {
		artistID, err := strconv.ParseUint(artist.StoreID, 10, 64)
		if err != nil {
			log.Errorf("can't parse uint64 from '%s'", artist.StoreID)
			continue
		}

		releases, err := albums.GetLatestArtistAlbums(f.Provider, artistID)
		if err != nil {
			if err == albums.ErrAlbumsNotFound {
				log.Debugf("Artist '%s' with id %s hasn't albums", artist.ArtistName, artist.StoreID)
				continue
			}

			log.Error(errors.Wrapf(err, "tried to get albums for '%s' with id %s", artist.ArtistName, artist.StoreID))
			continue
		}

		for _, release := range releases {
			err = db.DbMgr.EnsureReleaseExists(&db.Release{
				StoreName:  f.GetStoreName(),
				StoreID:    release.ID,
				ArtistName: artist.ArtistName,
				Title:      removeAlbumType(release.Attributes.Name),
				Poster:     release.Attributes.Artwork.GetLink(posterWidth, posterHeight),
				Released:   release.Attributes.ReleaseDate.Value,
			})
			if err != nil {
				log.Errorf("can't save release from '%s' with id '%s': %v", f.GetStoreName(), release.ID, err)
			}
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
}
