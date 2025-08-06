package models

import "gorm.io/gorm"

type PostModel struct {
	BaseModel
	User             *UserInfo    `json:"user"`
	UserId           uint         `json:"userId" form:"userId"`
	Content          string       `json:"content" form:"content" binding:"required" gorm:"size:1024"`
	Images           *StringArray `json:"images" form:"images" gorm:"type:json"`
	Topic            *TopicModel  `json:"topic" form:"topic" gorm:"-"`
	TopicId          *uint        `json:"topicId" form:"topicId"`
	Language         string       `json:"language" form:"language" gorm:"size:32"`
	Region           string       `json:"region" form:"region" gorm:"size:32"`
	LikeNum          uint         `json:"likeNum" form:"likeNum" gorm:"default:0"`
	CollectionNum    uint         `json:"collectionNum" form:"collectionNum" gorm:"default:0"`
	CommentNum       uint         `json:"commentNum" form:"commentNum" gorm:"default:0"`
	LikeStatus       uint         `json:"likeStatus" form:"likeStatus" gorm:"-"`
	CollectionStatus uint         `json:"collectionStatus" form:"collectionStatus" gorm:"-"`
}

func (post *PostModel) AfterFind(tx *gorm.DB) (err error) {
	if post.TopicId != nil {
		var topic TopicModel
		tx.Model(&TopicModel{}).Where("id = ?", *post.TopicId).First(&topic)
		post.Topic = &topic
	}
	return
}

// TopicModel topicStatus 0: 待审核 1: 审核通过 2: 审核未通过
// TopicType : 1.公共的，2：个人的
// /*
type TopicModel struct {
	BaseModel
	User        *UserInfo `json:"user"`
	UserId      uint      `json:"userId" form:"userId"`
	TopicType   uint      `json:"topicType" form:"topicType" binding:"required" gorm:"default:0"`
	Title       string    `json:"title" form:"title" gorm:"size:64"`
	Desc        string    `json:"desc" form:"desc" gorm:"size:256"`
	TitleKey    string    `json:"titleKey" form:"titleKey" gorm:"size:64"`
	DescKey     string    `json:"descKey" form:"descKey" gorm:"size:64"`
	TopicStatus uint      `json:"topicStatus" form:"topicStatus" gorm:"default:0"`
	Language    string    `json:"language" form:"language" gorm:"size:32"`
	Region      string    `json:"region" form:"region" gorm:"size:32"`
}
