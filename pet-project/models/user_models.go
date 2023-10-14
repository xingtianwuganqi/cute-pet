package models

import "gorm.io/gorm"

type UserInfo struct {
	gorm.Model
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Wx       string `json:"wx"`
	Location uint   `json:"location"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

type UserToken struct {
	gorm.Model
	UserId uint     `json:"userId"`
	User   UserInfo `json:"user" gorm:"foreignKey:UserId"`
	Token  string   `json:"token"`
}

func (UserToken) TableName() string {
	return "user_token"
}
