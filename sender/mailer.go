package sender

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

// MailConfig 邮箱服务器
type MailConfig struct {
	Host               string `required:"true"`
	Port               int    `required:"true"`
	Username           string `required:"true"`
	Password           string `required:"true"`
	Nickname           string
	InsecureSkipVerify bool
}

// New return mail server
func (mc MailConfig) New() *gomail.Dialer {
	dialer := gomail.NewPlainDialer(mc.Host, mc.Port, mc.Username, mc.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: mc.InsecureSkipVerify}
	dialer.LocalName = mc.Username
	return dialer
}
