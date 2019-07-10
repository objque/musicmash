package db

import "github.com/jinzhu/gorm"

var tables = []interface{}{
	&Release{},
	&LastAction{},
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

	fkeys := map[interface{}][][2]string{}

	for model, foreignKey := range fkeys {
		for _, fk := range foreignKey {
			if err := db.Debug().Model(model).AddForeignKey(
				fk[0], fk[1], "RESTRICT", "RESTRICT").Error; err != nil {
				return err
			}
		}
	}

	return db.Debug().Model(&Release{}).AddIndex("idx_created_at", "created_at").Error
}
