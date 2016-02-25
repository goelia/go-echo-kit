package models

import "github.com/jinzhu/gorm"

//AuthModel auth 公共结构体
type AuthModel struct {
	gorm.Model
	UserID uint `json:"user_id,omitempty"`
}
