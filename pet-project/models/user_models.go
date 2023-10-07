package models

import "gorm.io/gorm"

type UserInfo struct {
	gorm.Model
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Avatar    string    `json:"avatar"`
	Wx        string    `json:"wx"`
	UserToken UserToken `gorm:"foreignKey:UserId" json:"-"`
}

type UserToken struct {
	gorm.Model
	Token  string `json:"token"`
	UserId string `json:"userId"`
}
