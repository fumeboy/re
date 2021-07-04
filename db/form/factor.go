package form

import (
	"gorm.io/gorm"
	"re/db/client"
)

type Factor struct {
	gorm.Model
	Name        string
	Domain      int // Domain 和 Code 一起作为唯一标识
	Code        string
	Constructor string // 构造函数脚本；如果是空字符串，表示构造domain时的初始因子
}

type RelFactorOperation struct {
	OperationID, FactorID int
}

func init() {
	client.DB.AutoMigrate(&Factor{}, &RelFactorOperation{})
}
