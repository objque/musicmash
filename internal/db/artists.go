package db

import (
	"fmt"
	"html"
)

type Artist struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Poster     string `json:"poster"`
	Popularity int    `json:"popularity"`
	Followers  uint   `json:"followers"`
}

type ArtistMgr interface {
	EnsureArtistExists(artist *Artist) error
	GetAllArtists() ([]*Artist, error)
	SearchArtists(name string) ([]*Artist, error)
	GetArtist(id int64) (*Artist, error)
}

func (mgr *AppDatabaseMgr) GetAllArtists() ([]*Artist, error) {
	artists := []*Artist{}
	if err := mgr.db.Find(&artists).Error; err != nil {
		return nil, err
	}
	return artists, nil
}

func (mgr *AppDatabaseMgr) SearchArtists(name string) ([]*Artist, error) {
	artists := []*Artist{}
	name = fmt.Sprintf("%%%s%%", html.EscapeString(name))
	query := mgr.db.Where("name LIKE ?", name).Order("followers desc").Limit(100)
	if err := query.Find(&artists).Error; err != nil {
		return nil, err
	}
	return artists, nil
}

func (mgr *AppDatabaseMgr) EnsureArtistExists(artist *Artist) error {
	_, err := mgr.GetArtist(artist.ID)
	if err != nil {
		return mgr.db.Create(artist).Error
	}
	return nil
}

func (mgr *AppDatabaseMgr) GetArtist(id int64) (*Artist, error) {
	artist := Artist{}
	if err := mgr.db.First(&artist, id).Error; err != nil {
		return nil, err
	}
	return &artist, nil
}
