package db

import (
	"time"

	"github.com/pkg/errors"
)

type Feed struct {
	Date      time.Time
	Announced []*Release `json:"announced"`
	Released  []*Release `json:"released"`
}

type FeedMgr interface {
	GetUserFeedSince(userID string, date time.Time) (*Feed, error)
}

func (mgr *AppDatabaseMgr) GetUserFeedSince(userID string, since time.Time) (*Feed, error) {
	feed := &Feed{Date: since}
	var err error
	now := time.Now().UTC()
	feed.Released, err = mgr.GetReleasesForUserFilterByPeriod(userID, since, now)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get feed for user '%s'", userID)
	}

	feed.Announced, err = mgr.GetReleasesForUserSince(userID, time.Now().UTC())
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get future-feed for user '%s'", userID)
	}
	return feed, nil
}
