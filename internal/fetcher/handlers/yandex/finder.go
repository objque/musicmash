package yandex

import (
	"strings"

	"github.com/objque/musicmash/internal/clients/yandex"
	"github.com/objque/musicmash/internal/log"
	"github.com/pkg/errors"
)

func removeReleaseType(collection string) string {
	if strings.Contains(collection, " - Single") {
		log.Debugf("collectionName '%s' contains 'Single' that will be removed", collection)
		collection = strings.Replace(collection, " - Single", "", 1)
	}
	if strings.Contains(collection, " - EP") {
		log.Debugf("collectionName '%s' contains 'EP' that will be removed", collection)
		collection = strings.Replace(collection, " - EP", "", 1)
	}
	return strings.TrimSpace(collection)
}

func splitArtists(artists string) []string {
	// some releases have multiple artists
	// for example 1416189052 (iTunes id)
	// Wolfgang Muthspiel, Ambrose Akinmusire, Brad Mehldau, Larry Grenadier & Eric Harland - Where the River Goes
	if !strings.Contains(artists, " ") {
		return []string{artists}
	}

	// some releases with multiple artists have different separation: ['&', ',' ...]
	artists = strings.Replace(artists, ", ", " & ", -1)
	return strings.Split(artists, "&")
}

func searchArtistID(ya *yandex.Client, artistName string) (int, error) {
	log.Debugf("searching %s", artistName)
	res, err := ya.Search(artistName)
	if err != nil {
		return 0, errors.Wrapf(err, "tried to search artist %s", artistName)
	}

	log.Debugf("found %d artists", len(res.Artists.Items))
	for _, artist := range res.Artists.Items {
		if strings.ToLower(artist.Name) == strings.ToLower(artistName) {
			log.Debugf("100 artist name match (equals): %s artist_id: %d", artist.Name, artist.ID)
			return artist.ID, nil
		}
	}
	return 0, ArtistNotFoundErr
}

func find(ya *yandex.Client, releaseAuthor, releaseName string) (int, error) {
	releaseName = removeReleaseType(releaseName)
	for _, artistName := range splitArtists(releaseAuthor) {
		artistName = strings.TrimSpace(artistName)
		log.Debugf("searching artist %s in yandex music", artistName)

		artistID, err := searchArtistID(ya, artistName)
		if err != nil {
			if err == ArtistNotFoundErr {
				continue
			}

			log.Error(err)
			continue
		}

		albums, err := ya.GetArtistAlbums(artistID)
		if err != nil {
			return 0, errors.Wrapf(err, "tried to get artist (%d) albums", artistID)
		}

		for _, album := range albums {
			if strings.ToLower(releaseName) == strings.ToLower(album.Title) {
				log.Debugf("100 album name match: %s album_id: %d year: %d", album.Title, album.ID, album.ReleaseYear)
				return album.ID, nil
			}

			if strings.Contains(releaseName, album.Title) {
				// example cases:
				// 1. Body Count (Remix) [feat. Normani & Kehlani]
				// 2. Sun In Our Eyes (Don Diablo Remix)
				if strings.Contains(releaseName, album.Version) {
					log.Debugf("100 album name match: %s version: %s album_id: %d year: %d", album.Title, album.Version, album.ID, album.ReleaseYear)
					return album.ID, nil
				}

				// if releaseName is acoustic/mix/remix and yandexAlbum just release then we will have false-positive result :(
				log.Debugf("50 album name match: %s album_id: %d year: %d", album.Title, album.ID, album.ReleaseYear)
				return album.ID, nil
			}
		}
	}
	return 0, ReleaseNotFoundErr
}
