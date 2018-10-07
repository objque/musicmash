package db

import "github.com/jinzhu/gorm"

type Artist struct {
	Name string `gorm:"primary_key"`
}

type ArtistStoreInfo struct {
	id         uint64
	ArtistName string
	StoreName  string `gorm:"unique_index:idx_art_store_name_store_id"`
	StoreID    string `gorm:"unique_index:idx_art_store_name_store_id"`
}

type ArtistMgr interface {
	EnsureArtistExists(name string) error
}

func (mgr *AppDatabaseMgr) EnsureArtistExists(name string) error {
	info := ArtistStoreInfo{}
	err := mgr.db.Where("artist_name = ?", name).First(&info).Error
	if gorm.IsRecordNotFoundError(err) {
		return mgr.db.Create(Artist{Name: name}).Error
	}
	return err
}

type ArtistStoreInfoMgr interface {
	GetArtistsForStore(name string) ([]*ArtistStoreInfo, error)
	GetArtistFromStore(artistName, store string) ([]*ArtistStoreInfo, error)
	IsArtistExistsInStore(artistName, storeName, storeID string) bool
	EnsureArtistExistsInStore(artistName, storeName, storeID string) error
}

func (mgr *AppDatabaseMgr) GetArtistFromStore(name, store string) ([]*ArtistStoreInfo, error) {
	artists := []*ArtistStoreInfo{}
	if err := mgr.db.Where("artist_name = ? and store_name = ?", name, store).Find(&artists).Error; err != nil {
		return nil, err
	}
	return artists, nil
}

func (mgr *AppDatabaseMgr) GetArtistsForStore(name string) ([]*ArtistStoreInfo, error) {
	artists := []*ArtistStoreInfo{}
	if err := mgr.db.Where("store_name = ?", name).Find(&artists).Error; err != nil {
		return nil, err
	}
	return artists, nil
}

func (mgr *AppDatabaseMgr) IsArtistExistsInStore(artistName, storeName, storeID string) bool {
	info := ArtistStoreInfo{}
	err := mgr.db.Where("artist_name = ? and store_name = ? and store_id = ?", artistName, storeName, storeID).First(&info).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}
		return false
	}
	return true
}

func (mgr *AppDatabaseMgr) EnsureArtistExistsInStore(artistName, storeName, storeID string) error {
	if !mgr.IsArtistExistsInStore(artistName, storeName, storeID) {
		return mgr.db.Create(ArtistStoreInfo{ArtistName: artistName, StoreName: storeName, StoreID: storeID}).Error
	}
	return nil
}
