package musicvideos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/musicmash/musicmash/internal/clients/itunes"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/pkg/errors"
)

func getMusicVideos(provider *itunes.Provider, url string) (*Data, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.Token))
	provider.WaitRateLimit()
	resp, err := provider.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data := Data{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		b, readerErr := ioutil.ReadAll(decoder.Buffered())
		if readerErr != nil {
			return nil, fmt.Errorf("can't read all body because %v", readerErr)
		}
		return nil, errors.Wrapf(err, "looking json, but got %s", string(b))
	}

	return &data, nil
}

func GetArtistMusicVideos(provider *itunes.Provider, artistID uint64) ([]*MusicVideo, error) {
	log.Debugf("Getting music-videos by %d", artistID)
	url := fmt.Sprintf("%s/v1/catalog/us/artists/%v/music-videos?limit=100", provider.URL, artistID)
	data, err := getMusicVideos(provider, url)
	if err != nil {
		return nil, errors.Wrapf(err, "tried to get music-videos for %v", artistID)
	}

	musicVideos := data.MusicVideos
	for data.Next != "" {
		log.Debugf("Getting next page (%s) of music-videos for artist %v", data.Next, artistID)
		url := fmt.Sprintf("%s%s&limit=100", provider.URL, data.Next)
		data, err = getMusicVideos(provider, url)
		if err != nil {
			return nil, errors.Wrapf(err, "tried to get music-videos for %v", artistID)
		}
		musicVideos = append(musicVideos, data.MusicVideos...)
	}

	return musicVideos, nil
}

func isLatest(musicVideo *MusicVideo) bool {
	now := time.Now().UTC().Truncate(time.Hour * 24)
	yesterday := now.Add(-time.Hour * 48)
	return musicVideo.Attributes.ReleaseDate.Value.UTC().After(yesterday)
}

func GetLatestArtistMusicVideos(provider *itunes.Provider, artistID uint64) ([]*MusicVideo, error) {
	musicVideos, err := GetArtistMusicVideos(provider, artistID)
	if err != nil {
		return nil, err
	}

	latest := []*MusicVideo{}
	for _, release := range musicVideos {
		if !isLatest(release) {
			continue
		}

		latest = append(latest, release)
	}
	return latest, err
}
