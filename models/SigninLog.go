package models

//SigninLog 用户登录日志
type SigninLog struct {
	AuthModel
	IP         string //登录ip
	SigninType string //登录方式:code(校验码),local(用户+密码),oauth(第三方登录)
	AuthID     uint   //登录方式对应ID
}
