package models

type TopicModel struct {
	BaseModel
	User    *UserInfo `json:"user"`
	UserId  uint      `json:"userId" form:"userId"`
	Content string    `json:"content" form:"content" gorm:"size:1024"`
	Images  []string  `json:"images" form:"images" gorm:"type:json"`
	Tags    []uint    `json:"tags"`
}

type TagModel struct {
	BaseModel
	TagName string `json:"tagName" form:"tagName"`
}
