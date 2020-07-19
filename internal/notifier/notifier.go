package notifier

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/log"
)

const Action = "notify"

type Notifier struct {
	uri        *url.URL
	httpClient *http.Client
}

func New(rawurl string) (*Notifier, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	notifier := &Notifier{uri: uri, httpClient: newHTTPClient()}
	return notifier, nil
}

func (n *Notifier) Notify() error {
	// get last date when Action was successful
	last, err := db.Mgr.GetLastActionDate(Action)
	if err != nil {
		return fmt.Errorf("tried to get last_action for notify Action: %w", err)
	}

	// get releases since that date
	releases, err := db.Mgr.GetReleaseNotifications(last.Date)
	if err != nil {
		return fmt.Errorf("tried to get releases for notify, by got err: %w", err)
	}

	if len(releases) == 0 {
		log.Info("no new releases to notify")
		return nil
	}

	// group releases by user
	rels := groupReleases(releases)

	// send POST request with releases
	if err := n.sendReleases(rels); err != nil {
		return fmt.Errorf("tried to send releases, but got error from external service: %w", err)
	}

	return nil
}

func groupReleases(releases []*db.ReleaseNotification) []*Notification {
	group := map[string][]*db.InternalRelease{}
	for _, release := range releases {
		if _, ok := group[release.UserName]; !ok {
			group[release.UserName] = []*db.InternalRelease{}
		}

		rel := &db.InternalRelease{
			ArtistID:   release.ArtistID,
			ArtistName: release.ArtistName,
			Released:   release.Released,
			Poster:     release.Poster,
			Title:      release.Title,
			ItunesID:   release.ItunesID,
			SpotifyID:  release.SpotifyID,
			DeezerID:   release.DeezerID,
			Type:       release.Type,
			Explicit:   release.Explicit,
		}
		group[release.UserName] = append(group[release.UserName], rel)
	}

	notifications := make([]*Notification, len(group))
	for userName, releases := range group {
		notifications = append(notifications, &Notification{UserName: userName, Releases: releases})
	}

	return notifications
}
