package handlers

import (
	"time"

	"github.com/objque/musicmash/internal/itunes"
)

type StoreHandler interface {
	Fetch(releases []*itunes.LastRelease) []*ReleaseData
	NotifySubscribers(releases []*ReleaseData)
	GetStoreName() string
}

type ReleaseData struct {
	Date       time.Time
	ArtistName string
	StoreID    uint64
}
