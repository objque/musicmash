package apple

import (
	"sync"
	"time"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/musicmash/musicmash/internal/clients/itunes/artists"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/pkg/errors"
)

type Linker struct {
	jobs            chan []string
	mutex           sync.Mutex
	reservedArtists map[string]int
	client          *itunes.Provider
}

func NewLinker(url, token string) *Linker {
	return &Linker{
		jobs:            make(chan []string),
		mutex:           sync.Mutex{},
		reservedArtists: make(map[string]int),
		client:          itunes.NewProvider(url, token, time.Minute),
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
		dbArtists, err := db.DbMgr.GetArtistFromStore(artist, "itunes")
		if err != nil {
			log.Error(errors.Wrap(err, "tried to get artists while searching"))
			continue
		}

		if len(dbArtists) > 0 {
			continue
		}

		storeArtist, err := artists.SearchArtist(l.client, artist)
		if err != nil {
			if err == artists.ErrArtistNotFound {
				continue
			}
		}

		err = db.DbMgr.EnsureArtistExistsInStore(artist, "itunes", storeArtist.ID)
		if err != nil {
			log.Error(errors.Wrapf(err, "tried to save artist %s", storeArtist.ID))
		}
	}
}
