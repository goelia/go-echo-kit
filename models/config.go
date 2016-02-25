package models

import (
	"github.com/goelia/go-echo-kit/config"
	"github.com/goelia/go-echo-kit/sender"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
	ac AuthConfig
	mc sender.MailConfig
)

func init() {
	db = config.DB()
	ac = config.GetConfig().AuthConfig
}

// AuthConfig auth'config
type AuthConfig struct {
	AuthCodeRefreshExpSeconds int `default:"120"`
	AuthCodeExpSeconds        int `default:"300"`
	AuthCodeSigninTmpl        string
}
