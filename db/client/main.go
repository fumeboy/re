package client

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	db_, err := gorm.Open(sqlite.Open(""), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db_
}
