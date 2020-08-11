package user

import (
	"github.com/jinzhu/gorm"
	"github.com/yfedoruck/book/pkg/pg"
)

type Data struct {
	db       *gorm.DB
	Id       int
	Username string `json:"username" form:"username" query:"username"`
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
}

func New(db *pg.Postgres) *Data {
	d := new(Data)
	d.db = db.Get()
	return d
}

func (Data) TableName() string {
	return "account"
}

func (u *Data) Register() int {
	var lastInsertId int
	u.db.Create(u).Scan(&lastInsertId)

	return lastInsertId
}
