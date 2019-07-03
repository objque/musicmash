package fetcher

import (
	"sync"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/fetcher/services"
	"github.com/musicmash/musicmash/internal/fetcher/services/itunes"
	"github.com/musicmash/musicmash/internal/log"
)

func getServices() []services.Service {
	fetchers := []services.Service{}
	for name, store := range config.Config.Stores {
		// if fetching for current store is disabled
		if !store.Fetch {
			continue
		}

		switch name {
		case "itunes":
			fetchers = append(fetchers, itunes.NewService(store.URL, store.FetchWorkers, store.Meta["token"]))
		}
	}
	return fetchers
}

func fetchFromServices(services []services.Service) *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(len(services))

	// fetch from all services
	for i := range services {
		go services[i].FetchAndSave(&wg)
	}

	return &wg
}

func Fetch() {
	fetchFromServices(getServices()).Wait()

	// run callback
	log.Info("All stores were fetched")
}
