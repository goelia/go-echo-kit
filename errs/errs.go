package errs

import "fmt"

const (
// BadRequest 不合法的请求
	BadRequest = iota + 1000
// MethodNotAllowed 不允许的方法访问
	MethodNotAllowed
// NotValid 无效的
	NotValid
// NotSupported 不支持
	NotSupported
// Unauthorized 未授权
	Unauthorized
// AlreadyExists 已经存在
	AlreadyExists
// NotFound 不存在
	NotFound
// NotProvisioned 未配置
	NotProvisioned
// NotAssigned 未赋值
	NotAssigned
// NotExpired 未过期
	NotExpired
// Expired 已过期
	Expired
)

var errorText = map[int]string{
	BadRequest:"不合法的请求",
	MethodNotAllowed:"不允许的方法访问",
	NotValid :"无效的",
	NotSupported :"不支持",
	Unauthorized :"未授权",
	AlreadyExists :"已经存在",
	NotFound :"不存在",
	NotProvisioned:" 未配置",
	NotAssigned :"未赋值",
	NotExpired :"未过期",
	Expired :"已过期",
}

func ErrorText(code int) string {
	return errorText[code]
}

// Err response errors
type Err struct {
	Code        int     `json:"code,omitempty"`
	Message     string  `json:"message,omitempty"`
	Description string  `json:"description,omitempty"`
	Errors      []Err `json:"errors,omitempty"`
}

func (e *Err)Error() string {
	if e.Message == "" {
		e.Message = errorText[e.Code]
	}
	return fmt.Sprintf("%s", e.Message)
}

