package models

// 登录注册信息

type RegisterInfo struct {
	Phone    string `form:"phone" json:"phone"`
	Password string `form:"password" json:"password"  binding:"required"`
	Code     string `form:"code" json:"code"  binding:"required"`
	Email    string `form:"email" json:"email"`
}

type LoginInfo struct {
	Phone    string `form:"phone" json:"phone"`
	Password string `form:"password" json:"password"  binding:"required"`
	Email    string `form:"email" json:"email"`
}

type LoginUserInfo struct {
	ID     uint   `json:"id" form:"id"`
	Phone  string `json:"phone" form:"phone"`
	Email  string `json:"email" form:"email"`
	Avatar string `json:"avatar" form:"avatar"`
	Token  string `json:"token" form:"token"`
}

type EmptyModel struct {
}

// 宠物

type PetCustomActionCreateModel struct {
	BaseModel
	UserId     uint   `json:"userId" form:"userId"`
	CustomName string `json:"customName" form:"customName"`
	CustomIcon string `json:"customIcon" form:"customIcon"`
}

type PageModel struct {
	PageNum  int `json:"pageNum" form:"pageNum" binding:"required"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"required"`
}

type UploadPasswordModel struct {
	Password        string `json:"password" form:"password"`
	NewPassword     string `json:"newPassword" form:"newPassword"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword"`
}

type SendCodeModel struct {
	Phone string `json:"phone" form:"phone"`
	Email string `json:"email" form:"email"`
	Code  string `json:"code" form:"code" binding:"required"`
}

type UploadUserInfoModel struct {
	Avatar   string `json:"avatar" form:"avatar"`
	Username string `json:"username" form:"username"`
}

type RecordListModel struct {
	PageModel
	PetInfoId  uint  `json:"petInfoId" form:"petInfoId"`
	CategoryId *uint `json:"categoryId" form:"categoryId"`
}

// 帖子相关

type TopicStatusModel struct {
	TopicId uint `json:"topicId" form:"topicId"`
	Status  uint `json:"status" form:"status"`
}

type PostStatusModel struct {
	PostId uint `json:"postId" form:"postId"`
	Status uint `json:"status" form:"status"`
}

type MessageListType struct {
	PageModel
	MessageType uint `json:"messageType" form:"messageType"`
}

type CommentListModel struct {
	PageModel
	TopicId uint `json:"topicId" form:"topicId"`
}

type ReplyListModel struct {
	PageModel
	CommentId uint `json:"commentId" form:"commentId"`
}

type CategoryTypeModel struct {
	CategoryType uint `json:"categoryType" form:"categoryType"` //0:全部的 1:自己的
	PageModel
}
