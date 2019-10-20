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

func NewService(provider *itunes.Provider, fetchWorkers int) *Fetcher {
	return &Fetcher{
		Provider:     provider,
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

func (f *Fetcher) fetchWorker(id int, artists <-chan *db.ArtistStoreInfo, wg *sync.WaitGroup) {
	for artist := range artists {
		artistID, err := strconv.ParseUint(artist.StoreID, 10, 64)
		if err != nil {
			log.Errorf("can't parse uint64 from '%s'", artist.StoreID)
			wg.Done()
			continue
		}

		releases, err := albums.GetLatestArtistAlbums(f.Provider, artistID)
		if err != nil {
			if err == albums.ErrAlbumsNotFound {
				log.Debugf("artist with id %s hasn't albums", artist.StoreID)
				wg.Done()
				continue
			}

			log.Error(errors.Wrapf(err, "tried to get albums for artist with id %s", artist.StoreID))
			wg.Done()
			continue
		}

		for _, release := range releases {
			err = db.DbMgr.EnsureReleaseExists(&db.Release{
				StoreName: f.GetStoreName(),
				StoreID:   release.ID,
				ArtistID:  artist.ArtistID,
				Title:     removeAlbumType(release.Attributes.Name),
				Poster:    release.Attributes.Artwork.GetLink(posterWidth, posterHeight),
				Released:  release.Attributes.ReleaseDate.Value,
			})
			if err != nil {
				log.Errorf("can't save release from '%s' with id '%s': %v", f.GetStoreName(), release.ID, err)
			}
		}
		wg.Done()
	}
	log.Debugf("worker #%d finish fetching", id)
}

func (f *Fetcher) FetchAndSave(wg *sync.WaitGroup, storeArtists []*db.ArtistStoreInfo) {
	jobs := make(chan *db.ArtistStoreInfo, len(storeArtists))
	jobsWaitGroup := sync.WaitGroup{}
	jobsWaitGroup.Add(len(storeArtists))

	// Starts up X workers, initially blocked because there are no jobs yet.
	for w := 1; w <= f.FetchWorkers; w++ {
		go f.fetchWorker(w, jobs, &jobsWaitGroup)
	}

	// Here we send `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for _, artist := range storeArtists {
		jobs <- artist
	}

	jobsWaitGroup.Wait()
	close(jobs)
	wg.Done()
}
