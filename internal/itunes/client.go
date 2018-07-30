package itunes

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	releasedTag = `<time data-test-we-datetime datetime="`
	IDTag       = `class="featured-album targeted-link"`
)

func decode(buffer []byte) (*LastRelease, error) {
	parts := strings.Split(string(buffer), releasedTag)
	if len(parts) != 2 {
		return nil, errors.New("can't decode: after split by a time-html tag we have not 2 parts")
	}

	// Jul 18, 2018" aria-label="July 18 ...
	released := strings.Split(parts[1], `"`)[0]
	t, err := parseTime(released)
	if err != nil {
		return nil, errors.Wrapf(err, "can't parse time '%s'", released)
	}

	xxx := strings.Split(parts[0], IDTag)[0]
	urls := strings.Split(xxx, `<a href="`)
	return &LastRelease{
		Date: *t,
		URL:  strings.Replace(urls[len(urls)-1], `" `, "", 1),
	}, nil
}

func GetArtistInfo(url string) (*LastRelease, error) {
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
	return info, nil
}
