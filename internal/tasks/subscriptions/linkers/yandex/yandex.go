package yandex

import (
	"strconv"
	"sync"

	"github.com/musicmash/musicmash/internal/clients/yandex"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/pkg/errors"
)

type Linker struct {
	jobs            chan []string
	mutex           sync.Mutex
	reservedArtists map[string]int
	client          *yandex.Client
}

func NewLinker(url string) *Linker {
	return &Linker{
		jobs:            make(chan []string),
		mutex:           sync.Mutex{},
		reservedArtists: make(map[string]int),
		client:          yandex.New(url),
	}
}

func (l *Linker) reserveArtists(artists []string) []string {
	// reserve artists that current worker will search
	l.mutex.Lock()
	filteredArtists := []string{}
	for _, artist := range artists {
		if _, ok := l.reservedArtists[artist]; !ok {
			// no one doesn't search this artist right now
			l.reservedArtists[artist] = 1
			filteredArtists = append(filteredArtists, artist)
		} else {
			log.Debugln("someone already fetching artist", artist)
		}
	}
	l.mutex.Unlock()
	return filteredArtists
}

func (l *Linker) freeReservedArtists(artists []string) {
	l.mutex.Lock()
	for _, artist := range artists {
		delete(l.reservedArtists, artist)
	}
	l.mutex.Unlock()
}

func (l *Linker) SearchArtists(userArtists []string) {
	userArtists = l.reserveArtists(userArtists)
	defer func() { l.freeReservedArtists(userArtists) }()

	for _, artist := range userArtists {
		dbArtists, err := db.DbMgr.GetArtistFromStore(artist, "yandex")
		if err != nil {
			log.Error(errors.Wrap(err, "tried to get artists while searching"))
			continue
		}

		if len(dbArtists) > 0 {
			continue
		}

		storeArtists, err := l.client.Search(artist)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to search artist in the yandex: %s", artist))
			continue
		}

		if len(storeArtists.Artists.Items) == 0 {
			continue
		}

		storeArtist := storeArtists.Artists.Items[0]
		err = db.DbMgr.EnsureArtistExistsInStore(artist, "yandex", strconv.Itoa(storeArtist.ID))
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to save artist %d", storeArtist.ID))
		}
	}
}
