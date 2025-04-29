package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func Connect() (*gorm.DB, error) {
	dsn := "user=pg password=2280 port=5432 host=database dbname=connections sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return db, err
}
