package db

import sq "github.com/Masterminds/squirrel"

type Association struct {
	ID        int64  `json:"-"         db:"id"`
	ArtistID  int64  `json:"artist_id" db:"artist_id"`
	StoreName string `json:"name"      db:"store_name"`
	StoreID   string `json:"id"        db:"store_id"`
}

type AssociationOpts struct {
	ArtistID  int64
	StoreName string
	StoreID   string
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
	const query = "insert into associations (artist_id, store_name, store_id) values ($1, $2, $3)"

	_, err := mgr.newdb.Exec(query, artistID, storeName, storeID)

	return err
}

func applyAssociationsFilters(query sq.SelectBuilder, opts *AssociationOpts) sq.SelectBuilder {
	if opts.StoreID != "" {
		query = query.Where("store_id = ?", opts.StoreID)
	}
	if opts.StoreName != "" {
		query = query.Where("store_name = ?", opts.StoreName)
	}
	if opts.ArtistID != 0 {
		query = query.Where("artist_id = ?", opts.ArtistID)
	}
	return query
}

func (mgr *AppDatabaseMgr) FindAssociations(opts *AssociationOpts) ([]*Association, error) {
	query := sq.Select("artist_id", "store_name", "store_id").From("associations")
	if opts != nil {
		query = applyAssociationsFilters(query, opts)
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	associations := []*Association{}
	err = mgr.newdb.Select(&associations, sql, args...)
	if err != nil {
		return nil, err
	}

	return associations, nil
}
