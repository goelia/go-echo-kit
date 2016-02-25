package models

//CodeLog 校验码记录
type CodeLog struct {
	AuthModel
	Code       string
	CodeAuthID uint //对应CodeAuth的主键ID
}
