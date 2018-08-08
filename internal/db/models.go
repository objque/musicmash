package db

import (
	"github.com/jinzhu/gorm"
)

func NoOp(db *gorm.DB) error {
	return nil
}

var tables = []interface{}{
	&Artist{},
	&User{},
	&Release{},
	&LastFetch{},
	&Chat{},
	&Subscription{},
	&State{},
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
		Release{}: {
			{"artist_name", "artists(name)"},
		},
		Chat{}: {
			{"user_id", "users(id)"},
		},
		Subscription{}: {
			{"user_id", "users(id)"},
			{"artist_name", "artists(name)"},
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

	if err := db.Debug().Model(&Subscription{}).AddUniqueIndex(
		"idx_user_id_artist_name",
		"user_id", "artist_name").Error; err != nil {
		return err
	}

	return nil
}
