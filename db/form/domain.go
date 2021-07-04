package form

import (
	"gorm.io/gorm"
	"re/db/client"
)

type Domain struct {
	gorm.Model
	Name string
}

type DomainBridge struct { // domain -> domain
	From        int
	To          int
	Name        string
	Code        string
	Constructor string
}

func init() {
	client.DB.AutoMigrate(&DomainBridge{}, &Domain{})
}
