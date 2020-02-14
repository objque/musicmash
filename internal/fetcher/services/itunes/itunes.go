package itunes

import (
	"strconv"
	"strings"
	"sync"
	"time"

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
	SaveWorkers  int
}

func NewService(provider *itunes.Provider, fetchWorkers, saveWorkers int) *Fetcher {
	return &Fetcher{
		Provider:     provider,
		FetchWorkers: fetchWorkers,
		SaveWorkers:  saveWorkers,
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

type batch struct {
	ArtistID int64
	Releases []*albums.Album
}

//nolint:gocognit
func (f *Fetcher) fetchWorker(id int, artists <-chan *db.Association, releases chan<- *batch, wg *sync.WaitGroup) {
	log.Infof("Worker #%d is running", id)
	for artist := range artists {
		artistID, err := strconv.ParseUint(artist.StoreID, 10, 64)
		if err != nil {
			log.Errorf("can't parse uint64 from %s", artist.StoreID)
			wg.Done()
			continue
		}

		latestAlbums, err := albums.GetLatestArtistAlbums(f.Provider, artistID)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to get albums for artist with id %s", artist.StoreID))
			wg.Done()
			continue
		}

		if len(latestAlbums) == 0 {
			log.Debugf("Artist with id %s hasn't latest albums", artist.StoreID)
			wg.Done()
			continue
		}
		releases <- &batch{ArtistID: artist.ArtistID, Releases: latestAlbums}
	}
	log.Infof("Fetch worker #%d is finished", id)
}

func (f *Fetcher) saveWorker(id int, releases <-chan *batch, wg *sync.WaitGroup) {
	for batch := range releases {
		log.Debugf("Saving %d releases by %d", len(batch.Releases), batch.ArtistID)
		tx := db.DbMgr.Begin()
		now := time.Now().UTC()
		for _, release := range batch.Releases {
			title := removeAlbumType(release.Attributes.Name)
			err := tx.EnsureReleaseExists(&db.Release{
				CreatedAt: now,
				StoreName: f.GetStoreName(),
				StoreID:   release.ID,
				ArtistID:  batch.ArtistID,
				Title:     title,
				Poster:    release.Attributes.Artwork.GetLink(posterWidth, posterHeight),
				Released:  release.Attributes.ReleaseDate.Value,
			})
			if err != nil {
				log.Errorf("can't save release from %s with id %s: %v", f.GetStoreName(), release.ID, err)
			}
		}
		tx.Commit()
		wg.Done()
		log.Debugf("Finish saving releases by %d", batch.ArtistID)
	}
	log.Infof("Save worker #%d is finished", id)
}

func (f *Fetcher) FetchAndSave(wg *sync.WaitGroup, storeArtists []*db.Association) {
	jobs := make(chan *db.Association, len(storeArtists))
	jobsWaitGroup := sync.WaitGroup{}
	jobsWaitGroup.Add(len(storeArtists))
	releases := make(chan *batch, 250)

	for w := 1; w <= f.SaveWorkers; w++ {
		go f.saveWorker(w, releases, &jobsWaitGroup)
	}

	// Starts up X workers, initially blocked because there are no jobs yet.
	for w := 1; w <= f.FetchWorkers; w++ {
		go f.fetchWorker(w, jobs, releases, &jobsWaitGroup)
	}

	// Here we send `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for _, artist := range storeArtists {
		jobs <- artist
	}
	close(jobs)

	jobsWaitGroup.Wait()
	close(releases)
	wg.Done()
}
