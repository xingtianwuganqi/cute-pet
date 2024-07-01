package models

import "gorm.io/gorm"

// MessageModel 消息
/*
msg_type: # 1.点赞 2.收藏 3.评论 4.回复
msg_id: 帖子的id，可以是领养的也可以是秀宠的
reply_type: 1.评论 2.回复，3.关联的获取用户信息那张表的id
reply_id：评论的id
*/
type MessageModel struct {
	gorm.Model
	MessageType uint   `json:"message_type" form:"message_type" gorm:"not null,default:0"`
	MessageId   uint   `json:"message_id" form:"message_id" gorm:"not null default 0"`
	FromUid     uint   `json:"from_uid" form:"from_uid" gorm:"not null default 0"`
	ToUid       uint   `json:"to_uid" form:"to_uid" gorm:"not null default 0"`
	IsRead      bool   `json:"is_read" form:"is_read"`
	ReplyType   uint   `json:"reply_type" form:"reply_type" gorm:"not null default 0"`
	ReplyId     string `json:"reply_id" form:"reply_id" gorm:"not null default 0"`
}

type LikeMessageModel struct {
	gorm.Model
	LikeType   uint `json:"like_type" form:"like_type" form:"like_type" gorm:"not null,default:0"`
	LikeId     uint `json:"like_id" form:"like_id" gorm:"not null,default:0"`
	LikeStatus uint `json:"like_status" form:"like_status" gorm:"not null,default:0"`
	FromUid    uint `json:"from_uid" form:"from_uid" gorm:"not null,default:0"`
	ToUid      uint `json:"to_uid" form:"to_uid" gorm:"not null default 0"`
}

type CollectionMessageModel struct {
	gorm.Model
	CollectionType   uint `json:"collection_type" form:"collection_type" gorm:"not null default 0"`
	CollectionId     uint `json:"collection_id" form:"collection_id" gorm:"not null default 0"`
	CollectionStatus uint `json:"collection_status" form:"collection_status" gorm:"not null default 0"`
	FromUid          uint `json:"from_uid" form:"from_uid" gorm:"not null default 0"`
	ToUid            uint `json:"to_uid" form:"to_uid" gorm:"not null default 0"`
}

type CommentModel struct {
	gorm.Model
	TopicId   uint   `json:"topic_id" form:"topic_id" gorm:"not null default 0"`
	TopicType uint   `json:"topic_type" form:"topic_type" gorm:"not null default 0"`
	Content   string `json:"content" form:"content" gorm:"not null default '' size:256"`
	FromUid   uint   `json:"from_uid" form:"from_uid" gorm:"not null default 0"`
	ToUid     uint   `json:"to_uid" form:"to_uid" gorm:"not null default 0"`
}

type ReplayModel struct {
	gorm.Model
	CommentId uint   `json:"comment_id" form:"comment_id" gorm:"not null default 0"`
	ReplyId   uint   `json:"reply_id" form:"reply_id" gorm:"not null default 0"`
	Content   string `json:"content" form:"content" gorm:"not null default '' size:256"`
	FromUid   uint   `json:"from_uid" form:"FromUid" gorm:"not null default 0"`
	ToUid     uint   `json:"to_uid" form:"ToUid" gorm:"not null default 0"`
}
