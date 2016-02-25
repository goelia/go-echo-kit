package models

import (
	"github.com/goelia/go-echo-kit/config"
	"github.com/goelia/go-echo-kit/sender"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
	ac config.AuthConfig
	mc sender.MailConfig
)

func init() {
	db = config.DB()
	ac = config.GetConfig().AuthConfig
	mc = config.GetConfig().MailConfig
}

