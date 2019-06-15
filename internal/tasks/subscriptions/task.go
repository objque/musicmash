package subscriptions

import (
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/tasks/subscriptions/linkers/apple"
	"github.com/musicmash/musicmash/internal/tasks/subscriptions/linkers/yandex"
)

type job struct {
	UserName string
	Artists  []string
}

var (
	subscribeJobs = make(chan job)
	appleJobs     = make(chan []string)
	yandexJobs    = make(chan []string)
	appleLinker   *apple.Linker
	yandexLinker  *yandex.Linker
)

func InitWorkerPool() {
	for w := 1; w <= 3; w++ {
		go subscriber(subscribeJobs)
	}

	for name, store := range config.Config.Stores {
		switch name {
		case "itunes":
			appleLinker = apple.NewLinker(store.URL, store.Meta["token"])
			for w := 1; w <= 3; w++ {
				go runAppleWorker(appleJobs)
			}
		case "yandex":
			yandexLinker = yandex.NewLinker(store.URL)
			for w := 1; w <= 1; w++ {
				go runYandexWorker(yandexJobs)
			}
		}
	}
}

func SubscribeUserForArtists(userName string, artists []string) {
	subscribeJobs <- job{userName, artists}
}

func subscriber(jobs <-chan job) {
	// TODO (m.kalinin): handle concurrent requests from one user
	// TODO (m.kalinin): get artists that the user is not subscribed to
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}

			for _, artist := range job.Artists {
				_ = db.DbMgr.EnsureArtistExists(artist)
			}
			_ = db.DbMgr.SubscribeUserForArtists(job.UserName, job.Artists)
			linkArtists(job.Artists)
		}
	}
}

func linkArtists(artists []string) {
	appleJobs <- artists
	yandexJobs <- artists
}
