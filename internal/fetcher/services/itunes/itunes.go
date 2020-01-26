package itunes

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/musicmash/musicmash/internal/clients/itunes/albums"
	"github.com/musicmash/musicmash/internal/clients/itunes/songs"
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

type Release interface {
	GetID() string
	GetName() string
	GetPoster(width, height int) string
	GetReleaseDate() time.Time
}

type batch struct {
	ArtistID int64
	Type     string
	Releases []Release
}

//nolint:gocognit
func (f *Fetcher) fetchWorker(id int, artists <-chan *db.Association, releases chan<- *batch, done chan<- int) {
	log.Infof("Worker #%d is running", id)
	for association := range artists {
		artistStoreID, err := strconv.ParseUint(association.StoreID, 10, 64)
		if err != nil {
			log.Errorf("can't parse uint64 from %s", association.StoreID)
			continue
		}

		log.Debugf("Getting albums by artist %v associated with store id %v", association.ArtistID, artistStoreID)
		if rels := f.fetchAlbums(association.ArtistID, artistStoreID); rels != nil && len(rels) > 0 {
			releases <- &batch{ArtistID: association.ArtistID, Type: "album", Releases: rels}
		}

		log.Debugf("Getting songs by artist %v associated with store id %v", association.ArtistID, artistStoreID)
		if rels := f.fetchSongs(association.ArtistID, artistStoreID); rels != nil && len(rels) > 0 {
			releases <- &batch{ArtistID: association.ArtistID, Type: "song", Releases: rels}
		}
	}
	log.Infof("Fetch worker #%d is finished", id)
	done <- 1
}

func (f *Fetcher) fetchAlbums(artistID int64, storeID uint64) []Release {
	latestAlbums, err := albums.GetLatestArtistAlbums(f.Provider, storeID)
	if err != nil {
		log.Error(errors.Wrapf(err, "tried to get albums by artist %v associated with store id %v", artistID, storeID))
		return nil
	}

	rels := make([]Release, len(latestAlbums))
	for i := range latestAlbums {
		rels[i] = Release(latestAlbums[i])
	}
	return rels
}

func (f *Fetcher) fetchSongs(artistID int64, storeID uint64) []Release {
	latestSongs, err := songs.GetLatestArtistSongs(f.Provider, storeID)
	if err != nil {
		log.Error(errors.Wrapf(err, "tried to get songs by artist %v associated with store id %v", artistID, storeID))
		return nil
	}

	rels := []Release{}
	for i := range latestSongs {
		dbReleases, err := db.DbMgr.FindReleases(map[string]interface{}{
			"artist_id": artistID,
			"title":     removeAlbumType(latestSongs[i].GetAlbumName()),
		})
		if err != nil {
			log.Errorf("can't check if song from %s with id %s exists: %v",
				f.GetStoreName(), latestSongs[i].GetID(), err)
			continue
		}
		if len(dbReleases) == 0 {
			rels = append(rels, latestSongs[i])
		}
	}
	return rels
}

func (f *Fetcher) saveWorker(id int, releases <-chan *batch, done chan<- int) {
	for batch := range releases {
		log.Debugf("Saving %d releases by %d", len(batch.Releases), batch.ArtistID)
		tx := db.DbMgr.Begin()
		now := time.Now().UTC()
		for _, release := range batch.Releases {
			title := removeAlbumType(release.GetName())
			err := tx.EnsureReleaseExists(&db.Release{
				CreatedAt: now,
				StoreName: f.GetStoreName(),
				StoreID:   release.GetID(),
				ArtistID:  batch.ArtistID,
				Title:     title,
				Poster:    release.GetPoster(posterWidth, posterHeight),
				Released:  release.GetReleaseDate(),
				Type:      batch.Type,
			})
			if err != nil {
				log.Errorf("can't save release from %s with id %s: %v", f.GetStoreName(), release.GetID(), err)
			}
		}
		tx.Commit()
		log.Debugf("Finish saving releases by %d", batch.ArtistID)
	}
	log.Infof("Save worker #%d is finished", id)
	done <- 1
}

func (f *Fetcher) FetchAndSave(wg *sync.WaitGroup, storeArtists []*db.Association) {
	jobs := make(chan *db.Association, len(storeArtists))
	releases := make(chan *batch, 250)
	done := make(chan int, f.FetchWorkers+f.SaveWorkers)

	for w := 1; w <= f.SaveWorkers; w++ {
		go f.saveWorker(w, releases, done)
	}

	// Starts up X workers, initially blocked because there are no jobs yet.
	for w := 1; w <= f.FetchWorkers; w++ {
		go f.fetchWorker(w, jobs, releases, done)
	}

	// Here we send `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for _, artist := range storeArtists {
		jobs <- artist
	}
	close(jobs)

	for w := 1; w <= f.FetchWorkers; w++ {
		<-done
	}
	// indicate that's all releases we have.
	close(releases)

	// wait until save fetcher is done
	<-done
	wg.Done()
}
