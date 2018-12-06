package feed

import (
	"time"

	"github.com/musicmash/musicmash/internal/db"
)

func GetForUser(userName string, since, till time.Time) ([]*db.Release, error) {
	return db.DbMgr.GetReleasesForUserFilterByPeriod(userName, since, till)
}
