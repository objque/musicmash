package db

import (
	"fmt"
	"html"
	"time"
)

type Artist struct {
	ID         int64     `json:"id"          db:"id"`
	CreatedAt  time.Time `json:"-"           db:"created_at"`
	Name       string    `json:"name"        db:"name"`
	Poster     string    `json:"poster"      db:"poster"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
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
	const query = "select * from artists where name ilike $1 order by ID desc limit 100"

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
	const query = "insert into artists (name, poster, is_verified) values ($1, $2, $3) returning id"

	return mgr.newdb.QueryRow(query, artist.Name, artist.Poster, artist.IsVerified).Scan(&artist.ID)
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
