package handlers

import (
	"github.com/objque/musicmash/internal/db"
)

type StoreHandler interface {
	Fetch(releases []*db.Release)
	GetStoreName() string
}
