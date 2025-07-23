package models

type PostModel struct {
	BaseModel
	User    *UserInfo    `json:"user"`
	UserId  uint         `json:"userId" form:"userId"`
	Content string       `json:"content" form:"content" gorm:"size:1024"`
	Images  *StringArray `json:"images" form:"images" gorm:"type:json"`
	Topic   *TopicModel  `json:"topic" form:"topic" gorm:"-"`
	TopicId *uint        `json:"topicId" form:"topicId"`
}

type TopicModel struct {
	BaseModel
	User     *UserInfo `json:"user"`
	UserId   uint      `json:"userId" form:"userId"`
	Title    string    `json:"content" form:"content" gorm:"size:64"`
	Desc     string    `json:"desc" form:"desc" gorm:"size:256"`
	TitleKey string    `json:"titleKey" form:"titleKey" gorm:"size:64"`
	DescKey  string    `json:"descKey" form:"descKey" gorm:"size:64"`
}
