package dal

import (
	_ "re/db"
	"re/db/client"
	"re/db/form"
)

func DomainBridgeGet(from, to int) (form.DomainBridge, error) {
	var d form.DomainBridge
	sql := client.DB.Where(&form.DomainBridge{From: from, To: to}).First(&d)
	if sql.Error != nil {
		return form.DomainBridge{}, sql.Error
	}
	return d, nil
}

func DomainBridgeGet2(from_domain_id int, bridgeCode string) (form.DomainBridge, error) {
	var d form.DomainBridge
	sql := client.DB.Where(&form.DomainBridge{From: from_domain_id, Code: bridgeCode}).First(&d)
	if sql.Error != nil {
		return form.DomainBridge{}, sql.Error
	}
	return d, nil
}

func OperationGet(id int) (form.Operation, error) {
	var d form.Operation
	sql := client.DB.First(&d, id)
	if sql.Error != nil {
		return form.Operation{}, sql.Error
	}
	return d, nil
}

func EventDomainGet(eventID int) (form.RelEventDomain, error) {
	var d form.RelEventDomain
	sql := client.DB.Where(&form.RelEventDomain{EventID: eventID}).First(&d)
	if sql.Error != nil {
		return form.RelEventDomain{}, sql.Error
	}
	return d, nil
}
func FactorGet(domain_id int, code string) (form.Factor, error) {
	var d form.Factor
	sql := client.DB.Where(&form.Factor{Domain: domain_id, Code: code}).First(&d)
	if sql.Error != nil {
		return form.Factor{}, sql.Error
	}
	return d, nil
}
