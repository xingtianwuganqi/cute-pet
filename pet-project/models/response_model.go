package models

import "gorm.io/gorm"

type LoginInfo struct {
	Phone    string `form:"phone" json:"phone" binding:"required"`
	Password string `form:"password" json:"password"`
	Code     uint   `form:"code" json:"code"`
}

type LoginUserInfo struct {
	UserId uint   `json:"userId" form:"userId"`
	Phone  string `json:"phone" form:"phone"`
	Email  string `json:"email" form:"email"`
	Avatar string `json:"avatar" form:"avatar"`
	Token  string `json:"token" form:"token"`
}

type EmptyModel struct {
}

// PetCustomTypeInfo 自定义行为
type PetCustomTypeInfo struct {
	gorm.Model
	UserId     uint   `json:"userId" form:"userId"`
	CustomName string `json:"customName" form:"customName"`
	CustomIcon string `json:"customIcon" form:"customIcon"`
}
