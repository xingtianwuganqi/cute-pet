package models

import "gorm.io/gorm"

type LoginInfo struct {
	Phone    string `form:"phone" json:"phone" binding:"required"`
	Password string `form:"password" json:"password"`
	Code     uint   `form:"code" json:"code"`
}

type LoginUserInfo struct {
	UserId uint   `json:"userId"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Token  string `json:"token"`
}

type EmptyModel struct {
}

type PetCustomTypeInfo struct {
	gorm.Model
	UserId     uint   `json:"userId" form:"userId"`
	CustomName string `json:"customName" form:"customName"`
	CustomIcon string `json:"customIcon" form:"customIcon"`
}
