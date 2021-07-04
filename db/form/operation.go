package form

import (
	"gorm.io/gorm"
	"re/db/client"
)

type Operation struct {
	gorm.Model
	View   string
	Script string
}

func init() {
	client.DB.AutoMigrate(&Operation{})
}
