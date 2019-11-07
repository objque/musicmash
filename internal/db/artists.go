package db

import (
	"fmt"
	"html"
)

type Artist struct {
	ID         int64  `json:"id"        sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Name       string `json:"name"`
	Poster     string `json:"poster"`
	Popularity int    `json:"popularity"`
	Followers  uint   `json:"followers"`
}

type ArtistMgr interface {
	EnsureArtistExists(artist *Artist) error
	GetAllArtists() ([]*Artist, error)
	SearchArtists(name string) ([]*Artist, error)
	ValidateArtists(artists []int64) ([]int64, error)
	GetArtistWithFullInfo(id int64) (*Artist, error)
	GetArtistsWithFullInfo(ids []int64) ([]*Artist, error)
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
	return mgr.db.Create(artist).Error
}

func (mgr *AppDatabaseMgr) ValidateArtists(artists []int64) ([]int64, error) {
	result := []int64{}
	err := mgr.db.Table("artists").Where("id in (?)", artists).Pluck("id", &result).Error
	return result, err
}

func (mgr *AppDatabaseMgr) GetArtistWithFullInfo(id int64) (*Artist, error) {
	artist := Artist{}
	if err := mgr.db.First(&artist, id).Error; err != nil {
		return nil, err
	}
	return &artist, nil
}

func (mgr *AppDatabaseMgr) GetArtistsWithFullInfo(ids []int64) ([]*Artist, error) {
	artists := []*Artist{}
	if err := mgr.db.Where("id in (?)", ids).Find(&artists).Error; err != nil {
		return nil, err
	}
	return artists, nil
}
