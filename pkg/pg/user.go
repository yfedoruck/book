package pg

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/yfedoruck/book/pkg/user"
)

type User struct {
	user.Data
	db *gorm.DB
}

func (User) TableName() string {
	return "account"
}

func (u *User) Register() int {
	var lastInsertId int
	u.db.Create(u).Scan(&lastInsertId)

	return lastInsertId
}
