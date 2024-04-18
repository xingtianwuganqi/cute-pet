package models

import "gorm.io/gorm"

/*
location: 位置，默认0国内
*/

type UserInfo struct {
	gorm.Model
	Phone    string `json:"phone" form:"phone" gorm:"size:32"`
	Email    string `json:"email" form:"email" gorm:"size:32"`
	Username string `json:"username" form:"username" gorm:"size:32"`
	Password string `json:"password" form:"password" gorm:"size:64"`
	Avatar   string `json:"avatar" form:"avatar" gorm:"size:126"`
	Wx       string `json:"wx" form:"wx" gorm:"size:126"`
	Location uint   `json:"location" form:"location" gorm:"default:0"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

type SuggestionModel struct {
	gorm.Model
	UserId  uint `json:"userId"`
	User    UserInfo
	Contact string `json:"contact" gorm:"size:32"`
	Content string `json:"content" gorm:"size:256"`
}

func (SuggestionModel) TableName() string {
	return "suggestion"
}
