package services

import "sync"

type Service interface {
	FetchAndSave(done *sync.WaitGroup)
	ReFetchAndSave(done *sync.WaitGroup)
	GetStoreName() string
}
