package subscriptions

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/clients/itunes/v2"
	"github.com/objque/musicmash/internal/clients/itunes/v2/artists"
	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
	"github.com/objque/musicmash/internal/random"
	"github.com/pkg/errors"
)

const stateIDLength = 8

func findArtist(id int, provider *v2.Provider, jobs <-chan string, results chan<- string, done chan<- int) {
	for {
		userArtist, more := <-jobs
		if !more {
			break
		}

		dbArtist, err := db.DbMgr.FindArtistByName(userArtist)
		// artist already exists
		if err == nil {
			results <- dbArtist.Name
			continue
		}
		// another db err raised
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(errors.Wrapf(err, "tried to get artist '%s' from the db", userArtist))
			continue
		}

		artist, err := artists.SearchArtist(provider, userArtist)
		if err != nil {
			if err == artists.ArtistNotFoundErr {
				err = errors.Wrap(err, userArtist)
			}

			log.Error(err)
			continue
		}

		storeID, _ := strconv.ParseUint(artist.ID, 10, 64)
		err = db.DbMgr.CreateArtist(&db.Artist{Name: artist.Attributes.Name, StoreID: storeID})
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to add new artist '%s'", userArtist))
			continue
		}

		log.Debugf("found new artist '%s' storeID: %d", artist.Attributes.Name, storeID)
		results <- artist.Attributes.Name
	}
	done <- id
}

func subscribeUserForArtist(id int, userID string, jobs chan string, done chan int) {
	for {
		artistName, more := <-jobs
		if !more {
			log.Debugf("#%d subscribeUserForArtistWorker done", id)
			break
		}

		db.DbMgr.EnsureSubscriptionExists(&db.Subscription{ArtistName: artistName, UserID: userID})
		log.Debugf("subscribed user %s for %s", userID, artistName)
	}
	done <- id
}

func FindArtistsAndSubscribeUserTask(userID string, artists []string, provider *v2.Provider) (done chan bool, stateID string) {
	done = make(chan bool, 1)
	jobs := make(chan string, len(artists))
	results := make(chan string, len(artists))
	findWorkersDone := make(chan int, config.Config.Tasks.Subscriptions.FindArtistWorkers)
	subscribeWorkersDone := make(chan int, config.Config.Tasks.Subscriptions.SubscribeArtistWorkers)
	stateID = random.NewStringWithLength(stateIDLength)
	db.DbMgr.UpdateState(stateID, db.ProcessingState)
	startedAt := time.Now().UTC()

	for id := 1; id <= config.Config.Tasks.Subscriptions.FindArtistWorkers; id++ {
		go findArtist(id, provider, jobs, results, findWorkersDone)
	}

	for id := 1; id <= config.Config.Tasks.Subscriptions.SubscribeArtistWorkers; id++ {
		go subscribeUserForArtist(id, userID, results, subscribeWorkersDone)
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

		for id := 1; id <= config.Config.Tasks.Subscriptions.SubscribeArtistWorkers; id++ {
			log.Debugf("#%d subscribeArtistWorker done", <-subscribeWorkersDone)
		}

		err := db.DbMgr.UpdateState(stateID, db.CompleteState)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to update state '%s'", stateID))
		} else {
			log.Debugf("State '%s' was updated", stateID)
		}

		done <- true
		elapsed := time.Now().UTC().Sub(startedAt)
		log.Debugf("Finish fetch and subscribe user task. Elapsed time: %s", elapsed.String())
	}()
	return
}
