package db

import (
	"time"

	"github.com/pkg/errors"
)

type Feed struct {
	Date      time.Time  `json:"date"`
	Announced []*Release `json:"announced"`
	Released  []*Release `json:"released"`
}

type FeedMgr interface {
	GetUserFeedSince(userName string, date time.Time) (*Feed, error)
}

func (mgr *AppDatabaseMgr) GetUserFeedSince(userName string, since time.Time) (*Feed, error) {
	feed := &Feed{Date: since}
	var err error
	now := time.Now().UTC()
	feed.Released, err = mgr.GetReleasesForUserFilterByPeriod(userName, since, now)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get feed for user '%s'", userName)
	}

	feed.Announced, err = mgr.GetReleasesForUserSince(userName, time.Now().UTC())
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get future-feed for user '%s'", userName)
	}
	return feed, nil
}
