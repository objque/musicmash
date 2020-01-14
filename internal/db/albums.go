package db

import (
	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/log"
)

type Album struct {
	ID       uint64 `json:"id"`
	ArtistID int64  `json:"-"`
	Name     string `json:"name"`
}

type AlbumMgr interface {
	IsAlbumExists(album *Album) bool
	EnsureAlbumExists(album *Album) error
	GetAlbums(artistID int64) ([]*Album, error)
}

func (mgr *AppDatabaseMgr) IsAlbumExists(album *Album) bool {
	if err := mgr.db.Where("artist_id = ? and name = ?", album.ArtistID, album.Name).First(&album).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}

		log.Error(err)
		return false
	}
	return true
}

func (mgr *AppDatabaseMgr) EnsureAlbumExists(album *Album) error {
	if !mgr.IsAlbumExists(album) {
		return mgr.db.Create(album).Error
	}
	return nil
}

func (mgr *AppDatabaseMgr) GetAlbums(artistID int64) ([]*Album, error) {
	albums := []*Album{}
	err := mgr.db.Where("artist_id = ?", artistID).Find(&albums).Error
	if err != nil {
		return nil, err
	}
	return albums, nil
}
