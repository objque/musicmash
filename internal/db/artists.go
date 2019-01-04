package db

import (
	"encoding/json"
	"fmt"
	"html"

	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/log"
)

type Artist struct {
	Name   string `gorm:"primary_key" json:"name"`
	Poster string `json:"poster"`
}

type ArtistStoreInfo struct {
	ArtistName string `json:"-"`
	StoreName  string `gorm:"unique_index:idx_art_store_name_store_id" json:"name"`
	StoreID    string `gorm:"unique_index:idx_art_store_name_store_id" json:"id"`
}

func (r *ArtistStoreInfo) MarshalJSON() ([]byte, error) {
	s := struct {
		StoreURL  string `json:"url"`
		StoreName string `json:"name"`
		StoreID   string `json:"id"`
	}{
		StoreName: r.StoreName,
		StoreID:   r.StoreID,
	}

	if store, ok := config.Config.Stores[r.StoreName]; ok {
		s.StoreURL = fmt.Sprintf(store.ArtistURL, r.StoreID)
	} else {
		log.Warnf("artist_url for '%s' missed in config. User will receive empty link", r.StoreName)
	}
	return json.Marshal(&s)
}

type ArtistDetails struct {
	Artist
	Stores   []*ArtistStoreInfo     `gorm:"-" json:"stores"`
	Releases *ArtistDetailsReleases `json:"releases"`
}
type ArtistDetailsReleases struct {
	Announced []*Release `json:"announced"`
	Recent    []*Release `json:"released"`
}

type ArtistMgr interface {
	EnsureArtistExists(name string) error
	GetAllArtists() ([]*Artist, error)
	SearchArtists(name string) ([]*Artist, error)
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
	if err := mgr.db.Where("name LIKE ?", name).Order("name").Find(&artists).Error; err != nil {
		return nil, err
	}
	return artists, nil
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
	return err == nil
}

func (mgr *AppDatabaseMgr) EnsureArtistExistsInStore(artistName, storeName, storeID string) error {
	if !mgr.IsArtistExistsInStore(artistName, storeName, storeID) {
		return mgr.db.Create(ArtistStoreInfo{ArtistName: artistName, StoreName: storeName, StoreID: storeID}).Error
	}
	return nil
}

type ArtistDetailsMgr interface {
	GetArtistDetails(name string) (*ArtistDetails, error)
}

func (mgr *AppDatabaseMgr) GetArtistDetails(name string) (*ArtistDetails, error) {
	artist := Artist{}
	if err := mgr.db.Where("name = ?", name).First(&artist).Error; err != nil {
		return nil, err
	}

	stores := []*ArtistStoreInfo{}
	if err := mgr.db.Where("artist_name = ?", name).Find(&stores).Error; err != nil {
		return nil, err
	}

	announced, err := mgr.FindArtistAnnouncedReleases(name)
	if err != nil {
		return nil, err
	}

	released, err := mgr.FindArtistRecentReleases(name)
	if err != nil {
		return nil, err
	}

	details := &ArtistDetails{
		Artist: artist,
		Stores: stores,
		Releases: &ArtistDetailsReleases{
			Recent:    groupReleases(released),
			Announced: groupReleases(announced),
		},
	}
	return details, nil
}
