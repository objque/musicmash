package itunes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/log"
	"github.com/pkg/errors"
)

var (
	rxReleaseID   = regexp.MustCompile(`<a href.*?(\d+)([\?\/].*?)?" class="featured-album targeted-link"`)
	rxReleaseDate = regexp.MustCompile(`<time.*?>(.*?)<\/time>`)
)

func decode(buffer []byte) (*LastRelease, error) {
	body := string(buffer)
	// find release id
	releaseID := rxReleaseID.FindStringSubmatch(body)
	if len(releaseID) != 3 {
		return nil, errors.New("found too many substrings by id-regex")
	}
	id, err := strconv.ParseUint(releaseID[1], 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "can't parse uint from '%s'", releaseID[2])
	}

	// find release date
	released := rxReleaseDate.FindStringSubmatch(body)
	if len(released) != 2 {
		return nil, errors.New("found too many substrings by release-regex")
	}
	t, err := parseTime(released[1])
	if err != nil {
		return nil, errors.Wrapf(err, "can't parse time '%s'", released[1])
	}

	return &LastRelease{
		ID:   id,
		Date: *t,
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
