package pg

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	db       *gorm.DB
	Id       int
	Name     string
	Email    string
	Password string
}

func (u *User) Register() int {
	var lastInsertId int
	//err := u.db.QueryRow("INSERT into account (email,password,username) VALUES ($1,$2,$3) returning id;", u.Email, u.Password, u.Name).Scan(&lastInsertId)
	//fail.Check(err)

	return lastInsertId
}
