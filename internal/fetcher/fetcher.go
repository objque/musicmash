package fetcher

import (
	"sync"
	"time"

	itunes_client "github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/fetcher/services"
	"github.com/musicmash/musicmash/internal/fetcher/services/itunes"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/pkg/errors"
)

func getServices() []services.Service {
	fetchers := []services.Service{}
	for name, store := range config.Config.Stores {
		// if fetching for current store is disabled
		if !store.Fetch {
			continue
		}

		if name == "itunes" {
			itunesProvider := itunes_client.NewProvider(store.URL, store.Meta["token"], time.Minute)
			fetchers = append(fetchers, itunes.NewService(itunesProvider, store.FetchWorkers, store.SaveWorkers))
		}
	}
	return fetchers
}

func fetchFromServices(services []services.Service) *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(len(services))

	// fetch from all services
	for _, service := range services {
		storeArtists, err := db.DbMgr.GetAllAssociationsFromStore(service.GetStoreName())
		if err != nil {
			log.Error(errors.Wrapf(err, "can't receive artists from store: %s", service.GetStoreName()))
			wg.Done()
			continue
		}

		go service.FetchAndSave(&wg, storeArtists)
	}

	return &wg
}

func Fetch() {
	fetchFromServices(getServices()).Wait()

	// run callback
	log.Info("All stores were fetched")
}
