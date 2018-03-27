package orm

import (
	"database/sql"
)

type Orm interface {
	Insert(object interface{}) (int64, error)
	Select(object interface{}) error
}

func NewOrmr(sqlName string, params string) Orm {
	db, _ := sql.Open(sqlName, params)
	return &ormer{db: db}
}
