package services

import (
	"sync"

	"github.com/musicmash/musicmash/internal/db"
)

type Service interface {
	FetchAndSave(done *sync.WaitGroup, artists []*db.ArtistStoreInfo)
	GetStoreName() string
}
