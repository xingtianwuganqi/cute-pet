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
	Language string `json:"language" form:"language" gorm:"size:32"`
	Region   string `json:"region" form:"region" gorm:"size:32"`
}

type SuggestionModel struct {
	BaseModel
	User    *UserInfo `json:"user" form:"user"`
	UserId  uint      `json:"userId" form:"userId"`
	Contact string    `json:"contact" form:"contact" gorm:"size:32"`
	Content string    `json:"content" form:"content" gorm:"size:256"`
}

type IPInfo struct {
	IP          string `json:"ip"`
	Country     string `json:"country"`
	CountryName string `json:"country_name"`
	Region      string `json:"region"`
	City        string `json:"city"`
}
