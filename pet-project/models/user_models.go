package models

/*
location: 位置，默认0国内
*/

type UserInfo struct {
	BaseModel
	Phone    string `json:"phone" form:"phone" gorm:"size:32"`
	Email    string `json:"email" form:"email" gorm:"size:32"`
	Username string `json:"username" form:"username" gorm:"size:32"`
	Password string `json:"-" form:"-" gorm:"size:64"`
	Avatar   string `json:"avatar" form:"avatar" gorm:"size:126"`
	Wx       string `json:"wx" form:"wx" gorm:"size:126"`
	Location uint   `json:"location" form:"location" gorm:"default:0"`
}

type SuggestionModel struct {
	BaseModel
	UserId  uint   `json:"userId" gorm:"size:32"`
	Contact string `json:"contact" gorm:"size:32"`
	Content string `json:"content" gorm:"size:256"`
}
