package models

import (
	"github.com/goelia/go-echo-kit/errs"
	"github.com/jinzhu/gorm"
"github.com/goelia/go-echo-kit/utils"
)

//LocalAuth 系统登录
type LocalAuth struct {
	AuthModel
	Name     string `json:"name,omitempty" sql:"type:varchar(100);unique_index"` //手机/邮箱...
	Password []byte `json:"-"`
	Salt     []byte `json:"-"`
	Type     string //Name的类型
}

//Signin 系统登录
func (a *LocalAuth) Signin() error {
	pwd := a.Password
	if len(pwd) == 0 || len(a.Name) == 0 {
		return errs.New(errs.NotValid, "必填项不能为空")
	}
	err := db.Where("name = ?", a.Name).First(a).Error
	if err != nil {
		if err == gorm.RecordNotFound {
			return errs.New(errs.NotFound, "用户名不存在")
		}
		return err
	}
	if string(utils.Hash(a.Salt, pwd)) != string(a.Password) {
		return errs.New(errs.NotValid, "用户名或密码不正确")
	}

	return nil
}
