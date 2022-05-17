package model

import "gorm.io/gen"

type Method interface {
	// GetMaxVersionCount
	//
	// sql(select * from users order by version desc)
	GetMaxVersionCount() (gen.T, error)
}
