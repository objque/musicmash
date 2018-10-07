package subscribe

import (
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/tasks/subscribe/linkers/apple"
	"github.com/musicmash/musicmash/internal/tasks/subscribe/linkers/yandex"
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
		go subscriber(w, subscribeJobs)
	}

	for _, store := range config.Config.Stores {
		switch store.Name {
		case "itunes":
			appleLinker = apple.NewLinker(store.URL, store.Meta["token"])
			for w := 1; w <= 3; w++ {
				go runAppleWorker(w, appleJobs)
			}
		case "yandex":
			yandexLinker = yandex.NewLinker(store.URL)
			for w := 1; w <= 3; w++ {
				go runYandexWorker(w, yandexJobs)
			}
		}
	}
}

func SubscribeUserForArtists(userName string, artists []string) {
	subscribeJobs <- job{userName, artists}
}

func subscriber(worker int, jobs <-chan job) {
	// TODO (m.kalinin): handle concurrent requests from one user
	// TODO (m.kalinin): get artists that the user is not subscribed to
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}

			for _, artist := range job.Artists {
				db.DbMgr.EnsureArtistExists(artist)
			}
			db.DbMgr.SubscribeUserForArtists(job.UserName, job.Artists)
			linkArtists(job.Artists)
		}
	}
}

func linkArtists(artists []string) {
	appleJobs <- artists
	yandexJobs <- artists
}
