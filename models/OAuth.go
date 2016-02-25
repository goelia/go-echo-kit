package models

//OAuth 第三方登录
type OAuth struct {
	AuthModel
}

// TableName return OAuth's table name
func (OAuth) TableName() string {
	return "oauth"
}
