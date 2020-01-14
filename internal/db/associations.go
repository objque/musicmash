package db

type Association struct {
	ID        int64  `json:"-"`
	ArtistID  int64  `json:"artist_id"`
	StoreName string `json:"name"`
	StoreID   string `json:"id"`
}

type AssociationMgr interface {
	GetAllAssociationsFromStore(name string) ([]*Association, error)
	GetAssociationFromStore(artistID int64, store string) ([]*Association, error)
	IsAssociationExists(storeName, storeID string) bool
	EnsureAssociationExists(artistID int64, storeName, storeID string) error
}

func (mgr *AppDatabaseMgr) GetAssociationFromStore(artistID int64, store string) ([]*Association, error) {
	association := []*Association{}
	if err := mgr.db.Where("artist_id = ? and store_name = ?", artistID, store).Find(&association).Error; err != nil {
		return nil, err
	}
	return association, nil
}

func (mgr *AppDatabaseMgr) GetAllAssociationsFromStore(name string) ([]*Association, error) {
	association := []*Association{}
	if err := mgr.db.Where("store_name = ?", name).Find(&association).Error; err != nil {
		return nil, err
	}
	return association, nil
}

func (mgr *AppDatabaseMgr) IsAssociationExists(storeName, storeID string) bool {
	association := Association{}
	err := mgr.db.Where("store_name = ? and store_id = ?", storeName, storeID).First(&association).Error
	return err == nil
}

func (mgr *AppDatabaseMgr) EnsureAssociationExists(artistID int64, storeName, storeID string) error {
	if !mgr.IsAssociationExists(storeName, storeID) {
		return mgr.db.Create(&Association{ArtistID: artistID, StoreName: storeName, StoreID: storeID}).Error
	}
	return nil
}
