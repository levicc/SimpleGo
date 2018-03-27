package models

import "SimpleGo/simplego/orm"

type User struct {
	Id       int `PK`
	Username string
	Password string
}

var (
	ormr orm.Orm
)

func init() {
	ormr = orm.NewOrmr("mysql", "root:139145@tcp(localhost:3306)/gotest?charset=utf8")
}

func GetUserWithId(id int) (*User, error) {
	user := &User{Id: id}
	err := ormr.Select(user)
	return user, err
}
