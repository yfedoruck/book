package pg

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/yfedoruck/book/pkg/env"
	"github.com/yfedoruck/book/pkg/fail"
	"github.com/yfedoruck/book/pkg/lang"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Postgres struct {
	db *sql.DB
}

func (p *Postgres) Connect() {
	dbConf := Config()
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.Name)
	fmt.Println(dbInfo)
	var err error
	p.db, err = sql.Open("postgres", dbInfo)
	fail.Check(err)

	for i, connected := 0, false; connected == false && i < 4; i++ {
		err = p.db.Ping()
		if err == nil {
			connected = true
			return
		} else {
			log.Println("Error: Could not establish a connection with the database!", err, " but I still tried to connect...")
			time.Sleep(2 * time.Second)
		}
	}
	panic(err)
}

func (p *Postgres) Close() {
	err := p.db.Close()
	fail.Check(err)
}

func (p *Postgres) Tables() {
	files, err := filepath.Glob(env.BasePath() + "/sql/*.sql")
	fail.Check(err)

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		fail.Check(err)

		stmt, err := p.db.Prepare(string(data))
		fail.Check(err)

		_, err = stmt.Exec()
		fail.Check(err)
	}
}


func (p *Postgres) RegisterUser(username string, password string, email string) int {

	var lastInsertId int
	err := p.db.QueryRow("INSERT into account (email,password,username) VALUES ($1,$2,$3) returning id;", email, password, username).Scan(&lastInsertId)
	fail.Check(err)

	return lastInsertId
}

func (p *Postgres) LoginUser(username string) (int, string, error) {

	rows, err := p.db.Query("SELECT id, password FROM account WHERE username = $1 limit 1;", username)
	fail.Check(err)

	if rows.Next() == false {
		return 0, "", errors.New(lang.LoginErr)
	} else {
		var id int
		var password string
		err = rows.Scan(&id, &password)
		fail.Check(err)
		return id, password, nil
	}
}

func (p Postgres) IsUniqueUsername(username string) bool {

	// strings.Contains(err, "account_username_key")
	rows, err := p.db.Query("SELECT id FROM account WHERE username = $1 limit 1;", username)
	fail.Check(err)

	return rows.Next() == false
}

func (p Postgres) IsUniqueEmail(email string) bool {

	rows, err := p.db.Query("SELECT id FROM account WHERE email = $1 limit 1;", email)
	fail.Check(err)

	return rows.Next() == false
}

type Conf struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	Name     string `json:"Name"`
	Host     string `json:"Host"`
	Port     string `json:"Port"`
}

func Config() Conf {
	file, err := os.Open(env.BasePath() + filepath.FromSlash("/config/"+env.Domain()+"/postgres.json"))
	fail.Check(err)

	dbConf := Conf{}
	err = json.NewDecoder(file).Decode(&dbConf)
	fail.Check(err)

	return dbConf
}
