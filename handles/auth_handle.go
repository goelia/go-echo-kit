package handles

import (
	"github.com/goelia/go-echo-kit/conf"
	"github.com/goelia/go-echo-kit/errs"
	"github.com/goelia/go-echo-kit/models"
	"github.com/goelia/go-echo-kit/utils"
	"github.com/labstack/echo"
)

// Signin 用户登录
func Signin(c *echo.Context) error {
	var a struct {
		Name     string
		Password string
		Code     string
	}
	if err := c.Bind(&a); err != nil {
		return err
	}
	name := a.Name
	password := a.Password
	code := a.Code
	typ := utils.SigninType(name)
	println("type:", typ)
	switch typ {
	case "email", "mobile":
		switch {
		case len(password) > 0: //密码登录
		//la:=models.LocalAuth{}
		//if err := c.Bind(&la); err != nil {
		//	return err
		//}
		case len(code) > 0: //校验码登录
			ca := models.CodeAuth{}
			ca.Code = a.Code
			ca.Name = a.Name
			ip := c.Get("ip").(string)
			ca.IP = ip
			if err := ca.Signin(); err != nil {
				return err
			}
		default:
			return errs.New(errs.NotValid, "必填项不能为空")
		}

	default:
		return errs.New(errs.NotSupported, "不支持的登录方式")
	}

	//登录成功
	// 1.获取登录用户信息
	//e := models.Employee{}
	//employee, err := e.Fetch(a.UserID)
	//if err != nil {
	//	return err
	//}

	// 2.获取相关权限
	claims := map[string]interface{}{
		"sub":   "user.id",
		"name":  a.Name,
		"roles": "guest",
	}
	//if employee != nil {
	//	claims["employee"] = employee
	//}
	signingKey := conf.GetConfig().SigningKey
	jwtToken := utils.JwtToken(signingKey, claims)
	return c.JSON(200, map[string]interface{}{
		"token": jwtToken,
	})
}

// RefreshCode 发送验证码
func RefreshCode(c *echo.Context) error {
	a := models.CodeAuth{}
	if err := c.Bind(&a); err != nil {
		return err
	}
	seconds, err := a.RefreshCode("登录校验码", "简简单单")
	if err != nil {
		var res = map[string]interface{}{
			"code":    errs.NotValid,
			"message": err.Error(),
		}
		if seconds > 0 {
			res["seconds"] = seconds
		}
		return c.JSON(422, res)
	}
	return c.JSON(200, seconds)

}
