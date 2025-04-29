package migrations

import (
	"Golang_balancer/internal/config"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "29042025",
			Migrate: func(db *gorm.DB) error {
				return db.Migrator().CreateTable(&config.BucketDBConfig{})
			},
			Rollback: func(db *gorm.DB) error {
				return db.Migrator().DropTable("bucket_db_config")
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration success")
}
