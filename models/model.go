package models

import (
	"bitbucket.org/corn/conf"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

func init() {
	dbConfig := conf.GetConfig().DB
	db = dbConfig.New()
}

// Tables 将生成表的models集
func Tables() []interface{} {
	return []interface{}{
		&OAuth{},
		&LocalAuth{},
		&CodeAuth{},
		&CodeLog{},
		&User{},
		&SigninLog{},
	}
}

func CreateTables() {
	db.CreateTable(Tables())
}

func AutoMigrate() {
	for _, value := range Tables() {
		db.AutoMigrate(value)
	}
}
