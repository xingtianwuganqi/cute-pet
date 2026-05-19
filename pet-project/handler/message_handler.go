package handler

import (
	"pet-project/internal"
	"pet-project/models"
	"pet-project/response"
	"pet-project/service"
	"pet-project/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LikeMessageHandler(c *gin.Context) {
	logger.Logger.Info("LikeMessageHandler called", zap.String("clientIP", c.ClientIP()))

	userId := c.MustGet("userId").(uint)
	var statusModel models.LikeMessageModel
	if err := c.ShouldBind(&statusModel); err != nil {
		logger.Logger.Warn("Invalid parameters for LikeMessageHandler", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.LikeMessageService(userId, statusModel)
	if err != nil {
		logger.Logger.Error("Failed to process like message", zap.Error(err))
		if err.Error() == "query error" {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		} else {
			response.Fail(c, response.ApiCode.ServerErr, err.Error())
		}
		return
	}
	logger.Logger.Info("Successfully processed like message", zap.Uint("userId", userId))
	response.Success(c, result)
}

func CollectionMessageHandler(c *gin.Context) {
	logger.Logger.Info("CollectionMessageHandler called", zap.String("clientIP", c.ClientIP()))

	userId := c.MustGet("userId").(uint)
	var collectionModel models.CollectionMessageModel
	if err := c.ShouldBind(&collectionModel); err != nil {
		logger.Logger.Warn("Invalid parameters for CollectionMessageHandler", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.CollectionMessageService(userId, collectionModel)
	if err != nil {
		logger.Logger.Error("Failed to process collection message", zap.Error(err))
		if err.Error() == "query error" {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		} else {
			response.Fail(c, response.ApiCode.ServerErr, err.Error())
		}
		return
	}
	logger.Logger.Info("Successfully processed collection message", zap.Uint("userId", userId))
	response.Success(c, result)
}

// MessageListHandler 消息列表
func MessageListHandler(c *gin.Context) {
	logger.Logger.Info("MessageListHandler called", zap.String("clientIP", c.ClientIP()))

	var userId = c.MustGet("userId").(uint)
	var typeModel models.MessageListType
	if err := c.ShouldBind(&typeModel); err != nil {
		logger.Logger.Warn("Invalid parameters for MessageListHandler", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	msgList, err := service.MessageListService(userId, typeModel)
	if err != nil {
		logger.Logger.Error("Failed to get message list", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	logger.Logger.Info("Successfully retrieved message list", zap.Int("count", len(msgList)))
	response.Success(c, msgList)
}

// UnreadNumberHandler 未读消息数量
func UnreadNumberHandler(c *gin.Context) {
	logger.Logger.Info("UnreadNumberHandler called", zap.String("clientIP", c.ClientIP()))

	var userId = c.MustGet("userId").(uint)

	unreadNumbers, err := service.UnreadNumberService(userId)
	if err != nil {
		logger.Logger.Error("Failed to get unread numbers", zap.Error(err))
		response.Fail(c, response.ApiCode.ServerErr, err.Error())
		return
	}

	logger.Logger.Info("Successfully retrieved unread numbers", zap.Uint("userId", userId))
	response.Success(c, unreadNumbers)
}

// CommentHandler 评论
func CommentHandler(c *gin.Context) {
	logger.Logger.Info("CommentHandler called", zap.String("clientIP", c.ClientIP()))

	var commentModel models.CommentModel
	if err := c.ShouldBind(&commentModel); err != nil {
		logger.Logger.Warn("Invalid parameters for CommentHandler", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.CommentService(commentModel)
	if err != nil {
		logger.Logger.Error("Failed to process comment", zap.Error(err))
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	logger.Logger.Info("Successfully processed comment", zap.Uint("id", result.ID))
	response.Success(c, result)
}

func DeleteCommentHandler(c *gin.Context) {
	logger.Logger.Info("DeleteCommentHandler called", zap.String("clientIP", c.ClientIP()))

	commentId := internal.GetUintFromString(c.Param("commentId"))

	err := service.DeleteCommentService(commentId)
	if err != nil {
		logger.Logger.Error("Failed to delete comment", zap.Error(err), zap.Uint("commentId", commentId))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	logger.Logger.Info("Successfully deleted comment", zap.Uint("commentId", commentId))
	response.Success(c, nil)
}

// ReplyHandler 回复
func ReplyHandler(c *gin.Context) {
	logger.Logger.Info("ReplyHandler called", zap.String("clientIP", c.ClientIP()))

	var replyModel models.ReplyModel
	if err := c.ShouldBind(&replyModel); err != nil {
		logger.Logger.Warn("Invalid parameters for ReplyHandler", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.ReplyService(replyModel)
	if err != nil {
		logger.Logger.Error("Failed to process reply", zap.Error(err))
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	logger.Logger.Info("Successfully processed reply", zap.Uint("id", result.ID))
	response.Success(c, result)
}

func DeleteReplyHandler(c *gin.Context) {
	logger.Logger.Info("DeleteReplyHandler called", zap.String("clientIP", c.ClientIP()))

	replyId := internal.GetUintFromString(c.Param("replyId"))

	err := service.DeleteReplyService(replyId)
	if err != nil {
		logger.Logger.Error("Failed to delete reply", zap.Error(err), zap.Uint("replyId", replyId))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	logger.Logger.Info("Successfully deleted reply", zap.Uint("replyId", replyId))
	response.Success(c, nil)
}

func GetCommentListHandler(c *gin.Context) {
	logger.Logger.Info("GetCommentListHandler called", zap.String("clientIP", c.ClientIP()))

	var param models.CommentListModel
	if err := c.ShouldBindQuery(&param); err != nil {
		logger.Logger.Warn("Invalid parameters for GetCommentListHandler", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	commentList, err := service.GetCommentListService(param)
	if err != nil {
		logger.Logger.Error("Failed to get comment list", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	logger.Logger.Info("Successfully retrieved comment list", zap.Int("count", len(commentList)))
	response.Success(c, commentList)
}

func GetReplyListHandler(c *gin.Context) {
	logger.Logger.Info("GetReplyListHandler called", zap.String("clientIP", c.ClientIP()))

	var param models.ReplyListModel
	if err := c.ShouldBindQuery(&param); err != nil {
		logger.Logger.Warn("Invalid parameters for GetReplyListHandler", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	replyList, err := service.GetReplyListService(param)
	if err != nil {
		logger.Logger.Error("Failed to get reply list", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	logger.Logger.Info("Successfully retrieved reply list", zap.Int("count", len(replyList)))
	response.Success(c, replyList)
}

// 辅助函数，将字符串转换为uint
// func getUintFromString(s string) uint {
// 	var n uint
// 	_, _ = fmt.Sscanf(s, "%d", &n)
// 	return n
// }