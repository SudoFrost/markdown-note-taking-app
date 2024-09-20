package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("database.sqlite"), &gorm.Config{})
	return err
}

func init() {
	err := InitDB()
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&Category{})
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&Note{})
	if err != nil {
		panic(err)
	}
}
