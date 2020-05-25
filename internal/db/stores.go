package db

import (
	"database/sql"

	"github.com/musicmash/musicmash/internal/log"
)

type Store struct {
	Name string `json:"name"`
}

type StoreMgr interface {
	IsStoreExists(name string) bool
	EnsureStoreExists(name string) error
	CreateStore(name string) (*Store, error)
	GetAllStores() ([]*Store, error)
}

func (mgr *AppDatabaseMgr) IsStoreExists(name string) bool {
	const query = "select * from stores where name = $1"

	store := Store{}
	if err := mgr.newdb.Get(&store, query, name); err != nil {
		if err == sql.ErrNoRows {
			return false
		}

		log.Error(err)
		return false
	}
	return true
}

func (mgr *AppDatabaseMgr) EnsureStoreExists(name string) error {
	if mgr.IsStoreExists(name) {
		return nil
	}

	_, err := mgr.CreateStore(name)
	return err
}

func (mgr *AppDatabaseMgr) CreateStore(name string) (*Store, error) {
	const query = "insert into stores (name) values ($1)"

	_, err := mgr.newdb.Exec(query, name)

	store := &Store{Name: name}
	return store, err
}

func (mgr *AppDatabaseMgr) GetAllStores() ([]*Store, error) {
	const query = "select * from stores"

	stores := []*Store{}
	if err := mgr.newdb.Select(&stores, query); err != nil {
		return nil, err
	}

	return stores, nil
}
