package itunes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/log"
	"github.com/pkg/errors"
)

var (
	rxReleaseID          = regexp.MustCompile(`<a href.*\/\/(.*?\/){4}(\d+)([\?\/].*?)?" class="featured-album targeted-link"`)
	rxReleaseDate        = regexp.MustCompile(`<time.*?>(.*?)<\/time>`)
	rxComingReleaseDate  = regexp.MustCompile(`COMING (.*\d{4})`)
	htmlTagRelease       = []byte(`<h2 class="section__headline">Latest Release</h2>`)
	htmlTagComingRelease = []byte(`<h2 class="section__headline">Pre-Release</h2>`)
)

func findDate(body string, rx *regexp.Regexp, rxType string) (*time.Time, error) {
	released := rx.FindStringSubmatch(body)
	if len(released) != 2 {
		return nil, fmt.Errorf("found too many substrings by %s-regex: %v", rxType, released)
	}
	date, err := parseTime(released[1])
	if err != nil {
		return nil, errors.Wrapf(err, "can't parse datetime '%s'", released[1])
	}
	return date, nil
}

func findComingDate(body string) (*time.Time, error) {
	return findDate(body, rxComingReleaseDate, "coming")
}

func findReleaseDate(body string) (*time.Time, error) {
	return findDate(body, rxReleaseDate, "release")
}

func findReleaseID(body string) (*uint64, error) {
	releaseID := rxReleaseID.FindStringSubmatch(body)
	if len(releaseID) != 4 {
		return nil, fmt.Errorf("found too many substrings by id-regex: %v", releaseID)
	}

	id, err := strconv.ParseUint(releaseID[2], 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "can't parse uint from '%s'", releaseID[2])
	}
	return &id, err
}

func isComingRelease(buffer []byte) bool {
	return bytes.Contains(buffer, htmlTagComingRelease)
}

func isArtistInactive(buffer []byte) bool {
	return !bytes.Contains(buffer, htmlTagRelease) && !bytes.Contains(buffer, htmlTagComingRelease)
}

func decode(buffer []byte) (*LastRelease, error) {
	body := string(buffer)

	if isArtistInactive(buffer) {
		return nil, ArtistInactiveErr
	}

	releaseID, err := findReleaseID(body)
	if err != nil {
		return nil, err
	}

	findDate := findReleaseDate
	isComingRelease := isComingRelease(buffer)
	if isComingRelease {
		findDate = findComingDate
	}

	date, err := findDate(body)
	if err != nil {
		return nil, err
	}

	return &LastRelease{
		ID:       *releaseID,
		Date:     *date,
		IsComing: isComingRelease,
	}, nil
}

func GetArtistInfo(id uint64) (*LastRelease, error) {
	url := fmt.Sprintf("%s/%s/artist/%d", config.Config.Store.URL, config.Config.Store.Region, id)
	log.Debugf("Requesting '%s'...", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "can't receive page '%s'", url)
	}

	buffer, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, errors.Wrapf(err, "can't read response '%s'", url)
	}

	info, err := decode(buffer)
	if err != nil {
		// do not wrap inactive err
		// otherwise users can't compare received err with global InActiveErr error
		if err == ArtistInactiveErr {
			return nil, err
		}
		return nil, errors.Wrapf(err, "can't decode '%s'", url)
	}
	log.Debugf("Last release on '%s'", info.Date)
	return info, nil
}

func Lookup(id uint64) (*Release, error) {
	resp, err := http.Get(fmt.Sprintf("%s/lookup?id=%d", config.Config.Store.URL, id))
	if err != nil {
		return nil, err
	}

	searchResponse := SearchReleaseResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if searchResponse.Count == 0 {
		return nil, fmt.Errorf("release with given id '%d' not found", id)
	}

	if searchResponse.Count > 1 {
		return nil, fmt.Errorf("found more than one release with given id '%d'", id)
	}

	return searchResponse.Results[0], nil
}

func FindArtistID(artist string) (*Artist, error) {
	artist = strings.ToLower(artist)
	storeURL := fmt.Sprintf("%s/search?term=%s&media=music&limit=50", config.Config.Store.URL, url.QueryEscape(artist))
	resp, err := http.Get(storeURL)
	if err != nil {
		log.Error(errors.Wrapf(err, "tried to get '%s'", storeURL))
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("can't recieve answer from API. got status: %d for url: '%s'", resp.StatusCode, storeURL)
	}

	searchResponse := SearchReleaseResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, errors.Wrapf(err, "can't decode API response from '%s'", storeURL)
	}
	defer resp.Body.Close()

	if searchResponse.Count == 0 {
		return nil, ArtistNotFoundErr
	}

	for i, song := range searchResponse.Results {
		if artist == strings.ToLower(song.ArtistName) {
			log.Debugf("spent %d iterations to find storeID for '%s'", i, artist)
			artist := Artist{StoreID: uint64(song.ArtistID), Name: song.ArtistName}
			return &artist, nil
		}
	}
	return nil, errors.Wrapf(ArtistNotFoundErr, artist)
}
