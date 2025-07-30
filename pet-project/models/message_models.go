package models

// MessageModel 消息
/*
msg_type: # 1.点赞 2.收藏 3.评论 4.回复
msg_id: 帖子的id，可以是领养的也可以是秀宠的
reply_type: 1.评论 2.回复，3.关联的获取用户信息那张表的id
reply_id：评论的id
*/
type MessageModel struct {
	BaseModel
	MessageType uint `json:"messageType" form:"messageType" gorm:"default:0"`
	MessageId   uint `json:"messageId" form:"messageId" gorm:"default:0"`
	FromUid     uint `json:"fromUid" form:"fromUid" gorm:"default:0"`
	ToUid       uint `json:"toUid" form:"toUid" gorm:"default:0"`
	IsRead      bool `json:"isRead" form:"isRead" gorm:"default:false"`
	ReplyType   uint `json:"replyType" form:"replyType" gorm:"default:0"`
	ReplyId     uint `json:"replyId" form:"replyId" gorm:"default:0"`
}

// LikeMessageModel
/*
LikeType: # 1.点赞 2.收藏 3.评论 4.回复
*/
type LikeMessageModel struct {
	BaseModel
	LikeType   uint `json:"likeType" form:"likeType" form:"like_type" gorm:"default:0"`
	LikeId     uint `json:"likeId" form:"likeId" gorm:"default:0"`
	LikeStatus uint `json:"likeStatus" form:"likeStatus" gorm:"default:0"`
	FromUid    uint `json:"fromUid" form:"fromUid" gorm:"default:0"`
	ToUid      uint `json:"toUid" form:"toUid" gorm:"default:0"`
}

// CollectionMessageModel
/*
CollectionType  # 1.点赞 2.收藏 3.评论 4.回复
*/

type CollectionMessageModel struct {
	BaseModel
	CollectionType   uint `json:"collectionType" form:"collectionType" gorm:"default:0"`
	CollectionId     uint `json:"collectionId" form:"collectionId" gorm:"default:0"`
	CollectionStatus uint `json:"collectionStatus" form:"collectionStatus" gorm:"default:0"`
	FromUid          uint `json:"fromUid" form:"fromUid" gorm:"default:0"`
	ToUid            uint `json:"toUid" form:"toUid" gorm:"default:0"`
}

type CommentModel struct {
	BaseModel
	TopicId   uint   `json:"topicId" form:"topicId" gorm:"default:0"`
	TopicType uint   `json:"topicType" form:"topicType" gorm:"default:0"`
	Content   string `json:"content" form:"content" gorm:"size:256"`
	FromUid   uint   `json:"fromUid" form:"fromUid" gorm:"default:0"`
	ToUid     uint   `json:"toUid" form:"toUid" gorm:"default:0"`
}

// ReplayModel
// 通过commentId和replyId判断是回复的评论还是回复的回复
type ReplayModel struct {
	BaseModel
	CommentId uint   `json:"commentId" form:"commentId" gorm:"default:0"`
	ReplyId   uint   `json:"replyId" form:"replyId" gorm:"default:0"`
	Content   string `json:"content" form:"content" gorm:"size:256"`
	FromUid   uint   `json:"fromUid" form:"fromUid" gorm:"default:0"`
	ToUid     uint   `json:"toUid" form:"toUid" gorm:"default:0"`
}
