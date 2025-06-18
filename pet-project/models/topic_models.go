package models

type TopicModel struct {
	BaseModel
	User    UserInfo `json:"user"`
	UserId  uint     `json:"userId" form:"userId"`
	Content string   `json:"content" gorm:"size:1024"`
	Images  []string `json:"images" gorm:"size:256"`
	Tags    []uint   `json:"tags"`
}

type TagModel struct {
	BaseModel
	TagName string `json:"tagName" form:"tagName"`
}
