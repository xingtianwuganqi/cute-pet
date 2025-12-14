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

/*
topic_type:  领养1， 秀宠2，找宠3
*/
// CommentModel

type CommentModel struct {
	BaseModel
	TopicId   uint   `json:"topicId" form:"topicId" binding:"required" gorm:"default:0"`
	TopicType uint   `json:"topicType" form:"topicType" binding:"required" gorm:"default:0"`
	Content   string `json:"content" form:"content" binding:"required" gorm:"size:256"`
	FromUid   uint   `json:"fromUid" form:"fromUid" binding:"required" gorm:"default:0"`
	ToUid     uint   `json:"toUid" form:"toUid" binding:"required" gorm:"default:0"`
	// GORM 关联
	ReplyList []ReplyModel `json:"replyList" gorm:"foreignKey:CommentId;references:ID"` // 不加 form 标签，避免表单绑定
}

/*
	comment_id:表示在这条回复下
    reply_id: #表示回复目标的 id，如果 reply_type 是 comment 的话，那么 reply_id ＝ commit_id，如果 reply_type 是 reply 的话，这表示这条回复的父回复。
    reply_type: #表示回复的类型，因为回复可以是针对评论的回复（comment）值为1，也可以是针对回复的回复（reply），值为2， 通过这个字段来区分两种情景。
*/
// ReplayModel
// 通过commentId和replyId判断是回复的评论还是回复的回复
type ReplyModel struct {
	BaseModel
	CommentId uint   `json:"commentId" form:"commentId" binding:"required" gorm:"index"` // 建议加 index
	ReplyId   uint   `json:"replyId" form:"replyId" binding:"required" gorm:"default:0"`
	ReplyType uint   `json:"replyType" form:"replyType" binding:"required" gorm:"default:0"`
	Content   string `json:"content" form:"content" binding:"required" gorm:"size:256"`
	FromUid   uint   `json:"fromUid" form:"fromUid" binding:"required" gorm:"default:0"`
	ToUid     uint   `json:"toUid" form:"toUid" binding:"required" gorm:"default:0"`
}
