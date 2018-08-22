package handlers

import (
	"time"

	"github.com/objque/musicmash/internal/db"
)

type StoreHandler interface {
	Fetch(releases []*db.Release)
	NotifySubscribers(releases []*ReleaseData)
	GetStoreName() string
}

type ReleaseData struct {
	Date       time.Time
	ArtistName string
	StoreID    uint64
}
