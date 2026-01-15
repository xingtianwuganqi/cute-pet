package handler

import (
	"pet-project/models"
	"pet-project/response"
	"pet-project/service"

	"github.com/gin-gonic/gin"
)

func LikeMessageHandler(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var statusModel models.LikeMessageModel
	if err := c.ShouldBind(&statusModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.LikeMessageService(userId, statusModel)
	if err != nil {
		if err.Error() == "query error" {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		} else {
			response.Fail(c, response.ApiCode.ServerErr, err.Error())
		}
		return
	}
	response.Success(c, result)
}

func CollectionMessageHandler(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var collectionModel models.CollectionMessageModel
	if err := c.ShouldBind(&collectionModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.CollectionMessageService(userId, collectionModel)
	if err != nil {
		if err.Error() == "query error" {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		} else {
			response.Fail(c, response.ApiCode.ServerErr, err.Error())
		}
		return
	}
	response.Success(c, result)
}

// MessageListHandler 消息列表
func MessageListHandler(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var typeModel models.MessageListType
	if err := c.ShouldBind(&typeModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	msgList, err := service.MessageListService(userId, typeModel)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, msgList)
}

// UnreadNumberHandler 未读消息数量
func UnreadNumberHandler(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)

	unreadNumbers, err := service.UnreadNumberService(userId)
	if err != nil {
		response.Fail(c, response.ApiCode.ServerErr, err.Error())
		return
	}

	response.Success(c, unreadNumbers)
}

// CommentHandler 评论
func CommentHandler(c *gin.Context) {
	var commentModel models.CommentModel
	if err := c.ShouldBind(&commentModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.CommentService(commentModel)
	if err != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, result)
}

func DeleteCommentHandler(c *gin.Context) {
	commentId := getUintFromString(c.Param("commentId"))

	err := service.DeleteCommentService(commentId)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}

// ReplyHandler 回复
func ReplyHandler(c *gin.Context) {
	var replyModel models.ReplyModel
	if err := c.ShouldBind(&replyModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.ReplyService(replyModel)
	if err != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, result)
}

func DeleteReplyHandler(c *gin.Context) {
	replyId := getUintFromString(c.Param("replyId"))

	err := service.DeleteReplyService(replyId)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}

func GetCommentListHandler(c *gin.Context) {
	var param models.CommentListModel
	if err := c.ShouldBindQuery(&param); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	commentList, err := service.GetCommentListService(param)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, commentList)
}

func GetReplyListHandler(c *gin.Context) {
	var param models.ReplyListModel
	if err := c.ShouldBindQuery(&param); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	replyList, err := service.GetReplyListService(param)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, replyList)
}

// 辅助函数，将字符串转换为uint
// func getUintFromString(s string) uint {
// 	var n uint
// 	_, _ = fmt.Sscanf(s, "%d", &n)
// 	return n
// }