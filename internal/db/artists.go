package db

import (
	"fmt"
	"html"
)

type Artist struct {
	ID         int64  `json:"id"         db:"id"`
	Name       string `json:"name"       db:"name"`
	Poster     string `json:"poster"     db:"poster"`
	Popularity int    `json:"popularity" db:"popularity"`
	Followers  uint   `json:"followers"  db:"followers"`
}

func (mgr *AppDatabaseMgr) GetAllArtists() ([]*Artist, error) {
	const query = "select * from artists"

	artists := []*Artist{}
	err := mgr.newdb.Select(&artists, query)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (mgr *AppDatabaseMgr) SearchArtists(name string) ([]*Artist, error) {
	const query = "select * from artists where name like $1 order by followers desc limit 100"

	artists := []*Artist{}
	err := mgr.newdb.Select(&artists, query, fmt.Sprint("%", html.EscapeString(name), "%"))
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (mgr *AppDatabaseMgr) EnsureArtistExists(artist *Artist) error {
	_, err := mgr.GetArtist(artist.ID)
	if err != nil {
		return mgr.CreateArtist(artist)
	}

	return nil
}

func (mgr *AppDatabaseMgr) CreateArtist(artist *Artist) error {
	const query = "insert into artists (id, name, poster, popularity, followers) values ($1, $2, $3, $4, $5)"

	r, err := mgr.newdb.Exec(query, artist.ID, artist.Name, artist.Poster, artist.Popularity, artist.Followers)
	if err != nil {
		return err
	}

	id, _ := r.LastInsertId()
	artist.ID = id
	return nil
}

func (mgr *AppDatabaseMgr) GetArtist(id int64) (*Artist, error) {
	const query = "select * from artists where id = $1"

	artist := Artist{}
	err := mgr.newdb.Get(&artist, query, id)
	if err != nil {
		return nil, err
	}

	return &artist, nil
}
