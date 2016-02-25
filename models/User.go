package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//User 用户
type User struct {
	gorm.Model

	Username string `json:"username,omitempty"`                                  // validate:"required,excludesall=!@#$%^&*()_+-=:;?/0x2C" //登录用户名
	Nickname string `json:"nickname,omitempty"`                                  //花名(昵称)
	Email    string `json:"email,omitempty" sql:"unique_index" validate:"email"` //绑定的邮箱
	Mobile   string `json:"mobile,omitempty" sql:"unique_index"`                 //`validate:"mobile"`//绑定的手机号
	Name     string `json:"name,omitempty"`                                      //姓名
	//Age       uint8     `json:"age,omitempty" validate:"gte=0,lte=130"`
	Birthday  time.Time `json:"birthday,omitempty"`  //生日
	CreatedIP string    `json:"createdIP,omitempty"` //注册ip
	SigninAt  time.Time `json:"signinAt,omitempty"`
}

// TableName set User's table name
func (User) TableName() string {
	return "users"
}
