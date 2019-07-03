package services

import (
	"sync"

	"github.com/musicmash/artists/pkg/api/artists"
)

type Service interface {
	FetchAndSave(done *sync.WaitGroup, artists []*artists.StoreInfo)
	GetStoreName() string
}
