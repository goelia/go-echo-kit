package conf

import (
	"github.com/goelia/go-echo-kit/database"
	"github.com/goelia/go-echo-kit/sender"
	"github.com/koding/multiconfig"
)

var config Config

//Config 项目配置
type Config struct {
	SigningKey string
	Port       int
	Auth       AuthConfig
	DB         database.DBConfig
	Mailer     sender.MailConfig
}

type AuthConfig struct {
	AuthCodeRefreshExpSeconds int `default:"120"`
	AuthCodeExpSeconds        int `default:"300"`
	AuthCodeSigninTmpl        string
}

func init() {
	m := multiconfig.NewWithPath("config.toml")
	m.MustLoad(&config)
}

func GetConfig() Config {
	return config
}
