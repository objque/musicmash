package db

import (
	"sort"
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
	result := make(map[string]*Release)
	for _, value := range releases {
		// some releases might have equal titles, but from different artists
		key := strings.ToLower(value.ArtistName) + strings.ToLower(value.Title)
		if _, ok := result[key]; !ok {
			value.Stores = []*ReleaseStore{}
			result[key] = value
		}

		result[key].Stores = append(result[key].Stores, &ReleaseStore{
			StoreName: value.StoreName,
			StoreID:   value.StoreID,
		})

		// some releases haven't a poster, but if another
		// grouped release has a poster we should use it.
		if result[key].Poster == "" && value.Poster != "" {
			result[key].Poster = value.Poster
		}
	}

	releases = []*Release{}
	for _, release := range result {
		releases = append(releases, release)
	}

	sort.Slice(releases, func(i, j int) bool {
		return releases[i].Title < releases[j].Title
	})
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
