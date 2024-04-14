package models

import "gorm.io/gorm"

type TopicModel struct {
	gorm.Model
	User    UserInfo `json:"user" gorm:"foreignKey:UserId"`
	UserId  uint     `json:"userId" form:"userId"`
	Content string   `json:"content" gorm:"size:1024"`
	Images  []string `json:"images" gorm:"size:256"`
	Tags    []uint   `json:"tags"`
}

func (TopicModel) TableName() string {
	return "topic_models"
}

type TagModel struct {
	gorm.Model
	TagName string `json:"tagName" form:"tagName"`
}

func (TagModel) TableName() string {
	return "tag_models"
}
