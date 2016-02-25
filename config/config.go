package config

import (
	"github.com/goelia/go-echo-kit/database"
	"github.com/goelia/go-echo-kit/sender"
	"github.com/koding/multiconfig"
	"github.com/jinzhu/gorm"
	"github.com/goelia/go-echo-kit/models"
)

var (
	cfg Config
	db *gorm.DB
)

//Config 项目配置
type Config struct {
	SigningKey string
	Port       int
	AuthConfig models.AuthConfig
	DBConfig   database.DBConfig
	MailConfig sender.MailConfig
}

func init() {
	m := multiconfig.NewWithPath("config.toml")
	m.MustLoad(&cfg)

	db = cfg.DBConfig.New()
}

// GetConfig return config
func GetConfig() Config {
	return cfg
}

func DB() *gorm.DB {
	return db
}