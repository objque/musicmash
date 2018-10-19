package db

import (
	"strings"
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

func groupReleases(releases []*Release) []*Release {
	// key: lower(title), value: Release
	result := map[string]*Release{}
	for _, value := range releases {
		title := strings.ToLower(value.Title)
		if _, ok := result[title]; !ok {
			value.Stores = []*ReleaseStore{}
			result[title] = value
		}

		result[title].Stores = append(result[title].Stores, &ReleaseStore{
			value.StoreName,
			value.StoreID,
		})
	}

	releases = []*Release{}
	for _, release := range result {
		releases = append(releases, release)
	}
	return releases
}

func (mgr *AppDatabaseMgr) GetUserFeedSince(userName string, since time.Time) (*Feed, error) {
	var err error
	now := time.Now().UTC()
	released, err := mgr.GetReleasesForUserFilterByPeriod(userName, since, now)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get feed for user '%s'", userName)
	}

	future, err := mgr.GetReleasesForUserSince(userName, time.Now().UTC())
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get future-feed for user '%s'", userName)
	}

	feed := &Feed{
		Date:      since,
		Announced: groupReleases(future),
		Released:  groupReleases(released),
	}
	return feed, nil
}
