package notifier

import "github.com/musicmash/musicmash/internal/db"

type Notification struct {
	UserName string                `json:"user_name"`
	Releases []*db.InternalRelease `json:"releases"`
}
