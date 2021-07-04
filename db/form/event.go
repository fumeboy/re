package form

import (
	"re/db/client"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name string
}

type RelEventDomain struct {
	EventID, DomainID int
	Constructor       string
}

func init() {
	client.DB.AutoMigrate(&Event{}, &RelEventDomain{})
}
