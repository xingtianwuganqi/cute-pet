package models

import "gorm.io/gorm"

// MessageModel 消息
/*
message_type 消息类型
message
*/
type MessageModel struct {
	gorm.Model
	MessageType int    `json:"message_type" gorm:"not null,default:0"`
	Content     string `json:"content"`
}

func (MessageModel) TableName() string {
	return "messages"
}

type LikeMessageModel struct {
	gorm.Model
	LikeType int `json:"like_type" gorm:"not null,default:0"`
	LikeId   int `json:"like_id" gorm:"not null,default:0"`
}

func (LikeMessageModel) TableName() string {
	return "like_messages"
}

type CollectionMessageModel struct {
	gorm.Model
	CollectionType int `json:"collection_type" gorm:"not null default 0"`
	CollectionId   int `json:"collection_id" gorm:"not null default 0"`
}

func (CollectionMessageModel) TableName() string {
	return "collection_messages"
}
