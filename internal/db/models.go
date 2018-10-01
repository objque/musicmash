package db

import (
	"github.com/jinzhu/gorm"
)

func NoOp(db *gorm.DB) error {
	return nil
}

var tables = []interface{}{
	&Artist{},
	&ArtistStoreInfo{},
	&Store{},
	&User{},
	&Release{},
	&LastAction{},
	&Subscription{},
}

func CreateTables(db *gorm.DB) error {
	return db.AutoMigrate(tables...).Error
}

func DropAllTables(db *gorm.DB) error {
	return db.DropTable(tables...).Error
}

func CreateAll(db *gorm.DB) error {
	if err := CreateTables(db); err != nil {
		return err
	}

	fkeys := map[interface{}][][2]string{
		&ArtistStoreInfo{}: {
			{"artist_name", "artists(name)"},
			{"store_name", "stores(name)"},
		},
		Subscription{}: {
			{"user_name", "users(name)"},
			{"artist_name", "artists(name)"},
		},
		&Release{}: {
			{"artist_name", "artists(name)"},
			{"store_name", "stores(name)"},
		},
	}

	for model, model_fks := range fkeys {
		for _, fk := range model_fks {
			if err := db.Debug().Model(model).AddForeignKey(
				fk[0], fk[1], "RESTRICT", "RESTRICT").Error; err != nil {
				return err
			}
		}
	}

	//if err := db.Debug().Model(&Subscription{}).AddUniqueIndex(
	//	"idx_user_id_artist_name",
	//	"user_id", "artist_name").Error; err != nil {
	//	return err
	//}
	//
	//if err := db.Debug().Model(&Release{}).AddIndex(
	//	"idx_store_id", "store_id").Error; err != nil {
	//	return err
	//}
	//
	//if err := db.Debug().Model(&Store{}).AddIndex(
	//	"idx_store_type_release_id",
	//	"store_type", "release_id").Error; err != nil {
	//	return err
	//}

	return nil
}
