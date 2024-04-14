package models

type LoginInfo struct {
	Phone    string `form:"phone" json:"phone" binding:"required"`
	Password string `form:"password" json:"password"`
	Code     string `form:"code" json:"code"`
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
