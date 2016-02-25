package config

import (
	"github.com/goelia/go-echo-kit/database"
	"github.com/goelia/go-echo-kit/sender"
	"github.com/koding/multiconfig"
	"github.com/jinzhu/gorm"
)

var (
	cfg Config
	db *gorm.DB
)

//Config 项目配置
type Config struct {
	SigningKey string
	Port       int
	AuthConfig AuthConfig
	DBConfig   database.DBConfig
	MailConfig sender.MailConfig
}

// AuthConfig auth'config
type AuthConfig struct {
	AuthCodeRefreshExpSeconds int `default:"120"`
	AuthCodeExpSeconds        int `default:"300"`
	AuthCodeSigninTmpl        string
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