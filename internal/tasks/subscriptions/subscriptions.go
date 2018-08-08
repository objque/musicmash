package subscriptions

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/itunes"
	"github.com/objque/musicmash/internal/log"
	"github.com/objque/musicmash/internal/random"
	"github.com/pkg/errors"
)

const stateIDLength = 8

func findArtist(id int, jobs <-chan string, results chan<- uint64, done chan<- int) {
	for {
		userArtist, more := <-jobs
		if !more {
			done <- id
			break
		}

		dbArtist, err := db.DbMgr.FindArtistByName(userArtist)
		// artist already exists
		if err == nil {
			results <- dbArtist.StoreID
			continue
		}
		// another db err raised
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(errors.Wrapf(err, "tried to get artist '%s' from the db", userArtist))
			continue
		}

		artist, err := itunes.FindArtistID(userArtist)
		if err != nil {
			if err == itunes.ArtistNotFoundErr {
				err = errors.Wrap(err, userArtist)
			}

			log.Error(err)
			continue
		}

		err = db.DbMgr.CreateArtist(&db.Artist{Name: artist.Name, StoreID: artist.StoreID})
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to add new artist '%s'", userArtist))
			continue
		}

		log.Debugf("found new artist '%s' storeID: %d", artist.Name, artist.StoreID)
		results <- artist.StoreID
	}
}

func subscribeUserForArtist(id int, jobs chan uint64) {
	for {
		storeID, more := <-jobs
		if !more {
			log.Debugf("#%d subscribeUserForArtistWorker done", id)
			break
		}

		log.Debugf("subscribed user %s for %d", "objque", storeID)
	}
}

func FindArtistsAndSubscribeUserTask(userID string, artists []string) (done chan bool, stateID string) {
	done = make(chan bool, 1)
	jobs := make(chan string, len(artists))
	results := make(chan uint64, len(artists))
	findWorkersDone := make(chan int, config.Config.Tasks.Subscriptions.FindArtistWorkers)
	stateID = random.NewStringWithLength(stateIDLength)
	db.DbMgr.UpdateState(stateID, db.ProcessingState)
	startedAt := time.Now().UTC()

	for id := 1; id <= config.Config.Tasks.Subscriptions.FindArtistWorkers; id++ {
		go findArtist(id, jobs, results, findWorkersDone)
	}

	for id := 1; id <= config.Config.Tasks.Subscriptions.SubscribeArtistWorkers; id++ {
		go subscribeUserForArtist(id, results)
	}

	for _, artist := range artists {
		jobs <- artist
	}
	close(jobs)

	go func() {
		for id := 1; id <= config.Config.Tasks.Subscriptions.FindArtistWorkers; id++ {
			log.Debugf("#%d findArtistWorker done", <-findWorkersDone)
		}
		close(results)

		if err := db.DbMgr.UpdateState(stateID, db.CompleteState); err != nil {
			log.Error(errors.Wrapf(err, "tried to update state '%s'", stateID))
			return
		}

		elapsed := time.Now().UTC().Sub(startedAt)
		log.Debugf("State '%s' was updated", stateID)
		log.Debugf("Finish fetch and subscribe user task. Elapsed time: %s", elapsed.String())
		done <- true
	}()
	return
}
