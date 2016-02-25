package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
// mysql
	_ "github.com/go-sql-driver/mysql"
// postgres
	_ "github.com/lib/pq"
)

// DBConfig holds supported types by the multidc package
type DBConfig struct {
	Adapter       string `default:"mysql"`
	Host          string `default:"127.0.0.1"`
	Port          int    `default:"3306"`
	Name          string `required:"true"`
	Username      string `required:"true"`
	Password      string `required:"true"`
	SingularTable bool
	Debug         bool
}

// New return gorm.DB by dc
func (dc DBConfig) New() *gorm.DB {
	var (
		err error
		db gorm.DB

		adapter = strings.ToLower(dc.Adapter)
		url string
		username = dc.Username
		password = dc.Password
		host = dc.Host
		port = dc.Port
		dbName = dc.Name
	)
	// DBConfig := dc.DBConfig
	switch adapter {
	case "mysql":
		url = "%v:%v@tcp(%v:%d)/%v?charset=utf8&parseTime=True&loc=Local"
	case "postgres":
		// "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
		// or "user=%v password=%v dbname=%v sslmode=disable"
		url = "postgres://%v:%v@%v:%d/%v?sslmode=disable"
	default:
		log.Panicf("not support database adapter: %s", adapter)
	}

	db, err = gorm.Open(adapter, fmt.Sprintf(url, username, password, host, port, dbName))

	if err != nil {
		panic(err)
	}
	if dc.Debug {
		db.Debug()
	}
	db.LogMode(dc.Debug)
	db.SingularTable(dc.SingularTable)
	return &db
}
