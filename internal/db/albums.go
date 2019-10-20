package db

import (
	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/log"
)

type Album struct {
	ID       uint64 `json:"id"        gorm:"primary_key"   sql:"AUTO_INCREMENT"`
	ArtistID int64  `json:"artist_id" gorm:"unique_index:idx_album_art_id_name"`
	Name     string `json:"name"      gorm:"unique_index:idx_album_art_id_name"`
}

type AlbumMgr interface {
	IsAlbumExists(album *Album) bool
	EnsureAlbumExists(album *Album) error
}

func (mgr *AppDatabaseMgr) IsAlbumExists(album *Album) bool {
	if err := mgr.db.Where("name = ?", album.Name).First(&album).Error; err != nil {
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
