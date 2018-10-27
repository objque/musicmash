package yandex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	Session    *Session
	httpClient *http.Client
	URL        string
}

func New(url string) *Client {
	authURL := fmt.Sprintf("%s/api/v2.1/handlers/auth?external-domain=music.yandex.ru&overembed=no", url)
	req, _ := http.NewRequest(http.MethodGet, authURL, nil)
	req.Header.Set("X-Retpath-Y", "https%3A%2F%2Fmusic.yandex.ru")

	api := Client{httpClient: &http.Client{}, URL: url}
	resp, err := api.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&api.Session); err != nil {
		panic(err)
	}
	return &api
}

func (c *Client) do(url string, out interface{}) error {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Cookie", fmt.Sprintf("yandexuid=%s", c.Session.UID))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(&out)
}

func (c *Client) Search(text string) (*SearchResult, error) {
	searchURL := fmt.Sprintf("%s/handlers/music-search.jsx?type=all&text=%s", c.URL, url.QueryEscape(text))
	result := SearchResult{}
	if err := c.do(searchURL, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetArtistAlbums(id int) ([]*ArtistAlbum, error) {
	searchURL := fmt.Sprintf("%s/handlers/artist.jsx?what=albums&sort=year&artist=%d", c.URL, id)
	result := ArtistInfo{}
	if err := c.do(searchURL, &result); err != nil {
		return nil, err
	}
	return result.Albums, nil
}

func (c *Client) GetArtistLatestAlbum(artistID int) (*ArtistAlbum, error) {
	searchURL := fmt.Sprintf("%s/handlers/artist.jsx?what=albums&sort=year&artist=%d", c.URL, artistID)
	result := ArtistInfo{}
	if err := c.do(searchURL, &result); err != nil {
		return nil, err
	}

	if len(result.Albums) == 0 {
		return nil, ErrAlbumsNotFound
	}

	latest := result.Albums[0]
	for _, album := range result.Albums {
		if album.Released.Value.After(latest.Released.Value) {
			latest = album
		}
	}
	return latest, nil
}
