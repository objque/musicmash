package db

import "github.com/jinzhu/gorm"

type Association struct {
	ID        int64  `json:"-"`
	ArtistID  int64  `json:"artist_id"`
	StoreName string `json:"name"`
	StoreID   string `json:"id"`
}

type AssociationOpts struct {
	ArtistID  int64
	StoreName string
	StoreID   string
}

type AssociationMgr interface {
	GetAllAssociationsFromStore(name string) ([]*Association, error)
	GetAssociationFromStore(artistID int64, store string) ([]*Association, error)
	IsAssociationExists(storeName, storeID string) bool
	EnsureAssociationExists(artistID int64, storeName, storeID string) error
	FindAssociations(opts *AssociationOpts) ([]*Association, error)
}

func (mgr *AppDatabaseMgr) GetAssociationFromStore(artistID int64, store string) ([]*Association, error) {
	return mgr.FindAssociations(&AssociationOpts{ArtistID: artistID, StoreName: store})
}

func (mgr *AppDatabaseMgr) GetAllAssociationsFromStore(name string) ([]*Association, error) {
	return mgr.FindAssociations(&AssociationOpts{StoreName: name})
}

func (mgr *AppDatabaseMgr) IsAssociationExists(storeName, storeID string) bool {
	associations, _ := mgr.FindAssociations(&AssociationOpts{StoreName: storeName, StoreID: storeID})
	return len(associations) > 0
}

func (mgr *AppDatabaseMgr) EnsureAssociationExists(artistID int64, storeName, storeID string) error {
	return mgr.db.Create(&Association{ArtistID: artistID, StoreName: storeName, StoreID: storeID}).Error
}

func applyAssociationsFilters(db *gorm.DB, opts *AssociationOpts) *gorm.DB {
	if opts.StoreID != "" {
		db = db.Where("store_id = ?", opts.StoreID)
	}
	if opts.StoreName != "" {
		db = db.Where("store_name = ?", opts.StoreName)
	}
	if opts.ArtistID != 0 {
		db = db.Where("artist_id = ?", opts.ArtistID)
	}
	return db
}

func (mgr *AppDatabaseMgr) FindAssociations(opts *AssociationOpts) ([]*Association, error) {
	association := []*Association{}
	db := mgr.db
	if opts != nil {
		db = applyAssociationsFilters(db, opts)
	}
	if err := db.Find(&association).Error; err != nil {
		return nil, err
	}
	return association, nil
}
