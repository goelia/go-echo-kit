package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/goelia/go-echo-kit/errs"
	"github.com/goelia/go-echo-kit/utils"
	"github.com/jinzhu/gorm"
	"gopkg.in/gomail.v2"
)

//CodeAuth 校验码登录
type CodeAuth struct {
	AuthModel
	Name          string    `json:"name,omitempty" sql:"type:varchar(100);unique_index"` //手机/邮箱/微信...
	Code          string    `json:"code,omitempty"`                                      //校验码
	CodingAt      time.Time `json:"coding_at,omitempty"`                                 //当前校验码生成时间
	ExpireSeconds int       `json:"expireSeconds,omitempty"`                             //过期时间(单位:秒)
	Type          string    //Name的类型
	IP            string    `sql:"-"`
}

//Signin 校验码登录
func (a *CodeAuth) Signin() error {
	//1. 验证用户名格式
	if utils.SigninType(a.Name) == "" {
		return errs.New(errs.NotSupported, "不支持的登录方式")
	}
	//2. 查找数据库记录
	if err := db.Where("name=? and code=?", a.Name, a.Code).First(a).Error; err != nil {
		if err == gorm.RecordNotFound {
			return errs.New(errs.NotFound, "校验码不正确")
		}
		return err
	}
	//3. 校验校验码是否过期
	if a.UpdatedAt.Add(time.Second * time.Duration(a.ExpireSeconds)).Before(time.Now()) {
		//判断校验码是否过期
		return errs.New(errs.Expired, "校验码已过期")
	}
	//4. 校验通过, 更新校验码为过期.写入登录日志.更新User的修改日期为最新登录日期
	tx := db.Begin()
	//过期校验码
	if err := db.Model(a).Update("expire_seconds", 0).Error; err != nil {
		tx.Rollback()
		return err
	}
	//写入登录日志
	sl := SigninLog{
		IP:         a.IP,
		SigninType: "code",
		AuthID:     a.ID,
	}
	sl.UserID = a.UserID
	if err := db.Create(&sl).Error; err != nil {
		tx.Rollback()
		return err
	}

	//更新User最后修改时间
	u := User{}
	u.ID = a.ID
	if err := db.Model(&u).Update("signin_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// RefreshCode 发送校验码
// seconds 返回重试时间
func (a *CodeAuth) RefreshCode(subject, nickname string) (int, error) {
	st := utils.SigninType(a.Name)
	//1. 验证用户名格式
	if st == "" {
		return 0, errs.New(errs.NotSupported, "不支持的登录方式")
	}
	//2. 用户是否存在
	err := db.Where("name = ?", a.Name).First(a).Error

	tx := db.Begin()

	if err == nil {
		//存在用户, 校验重发间隔时间, 过期则刷新校验码
		exp := a.CodingAt.Add(time.Duration(ac.AuthCodeRefreshExpSeconds) * time.Second)
		sub := exp.Sub(time.Now()).Seconds()
		if int(sub) > 0 {
			tx.Rollback()
			return int(sub), errs.New(errs.NotExpired, fmt.Sprintf("请%d秒后重试", int(sub)))
		}
	} else if err == gorm.RecordNotFound {
		//不存在则注册为新用户,发送用户校验码, 新建校验码日志

		//新建用户
		user := User{
			Email: a.Name,
		}
		err = tx.Create(&user).Error
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		a.UserID = user.ID

	} else {
		return 0, err
	}
	code := strconv.Itoa(utils.RandNum())
	//刷新校验码
	a.Code = code
	a.CodingAt = time.Now()
	a.ExpireSeconds = ac.AuthCodeExpSeconds
	if err = tx.Save(a).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	//写入校验码日志
	cl := CodeLog{
		Code:       code,
		CodeAuthID: a.ID,
	}
	cl.UserID = a.UserID
	err = tx.Create(&cl).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	//发送校验码
	switch st {
	case "email":
		content := fmt.Sprintf(ac.AuthCodeSigninTmpl, nickname, code, ac.AuthCodeExpSeconds/60)
		m := gomail.NewMessage()
		dialer := mc.New()
		m.SetAddressHeader("From", dialer.LocalName, nickname)
		m.SetHeader("To", a.Name)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", content)

		if err = dialer.DialAndSend(m); err == nil {
			tx.Rollback()
			return 0, err
		}
		m.Reset()
	case "mobile":
		tx.Rollback()
		return 0, errs.New(errs.NotSupported, "不支持的登录方式")
	}

	tx.Commit()
	return ac.AuthCodeExpSeconds, nil
}
