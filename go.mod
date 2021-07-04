module re

go 1.16

require (
	github.com/d5/tengo/v2 v2.7.0
	github.com/pkg/errors v0.9.1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.8
)

replace errors => github.com/pkg/errors v0.9.1
