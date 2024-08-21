package models

import "gorm.io/gorm"

type RegisterInfo struct {
	Phone    string `form:"phone" json:"phone"`
	Password string `form:"password" json:"password"  binding:"required"`
	Code     int    `form:"code" json:"code"  binding:"required"`
	Email    string `form:"email" json:"email"`
}

type LoginInfo struct {
	Phone    string `form:"phone" json:"phone"`
	Password string `form:"password" json:"password"  binding:"required"`
	Email    string `form:"email" json:"email"`
}

type LoginUserInfo struct {
	ID     uint   `json:"ID" form:"ID"`
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

type PageModel struct {
	PageNum  int `json:"pageNum" form:"pageNum" binding:"required"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"required"`
}
