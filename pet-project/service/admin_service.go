package service

import (
	"pet-project/db"
	"pet-project/models"
	"pet-project/logger"

	"gorm.io/gorm/clause"
	"go.uber.org/zap"
)

// GetCommonCategoriesService 获取宠物分类的业务逻辑
func GetCommonCategoriesService() ([]models.RecordCategory, error) {
	logger.Logger.Info("GetCommonCategoriesService called")

	var petActionList []models.RecordCategory
	result := db.DB.Model(&models.RecordCategory{}).
		Where("user_id IS NULL").
		Find(&petActionList)
	if result.Error != nil {
		logger.Logger.Error("Database error in GetCommonCategoriesService", zap.Error(result.Error))
		return nil, result.Error
	}
	logger.Logger.Info("Successfully retrieved common categories", zap.Int("count", len(petActionList)))
	return petActionList, nil
}

// CreateCommonCategoryService 创建宠物分类的业务逻辑
func CreateCommonCategoryService(recordCategory models.RecordCategory) error {
	logger.Logger.Info("CreateCommonCategoryService called", zap.String("name", recordCategory.Name))

	recordCategory.UserId = nil
	result := db.DB.Omit(clause.Associations).Create(&recordCategory)
	if result.Error != nil {
		logger.Logger.Error("Database error in CreateCommonCategoryService", zap.Error(result.Error))
		return result.Error
	}
	logger.Logger.Info("Successfully created common category", zap.Uint("id", recordCategory.ID))
	return nil
}

// CreateCommonCategoryListService 批量添加宠物行为的业务逻辑
func CreateCommonCategoryListService(categories []models.RecordCategory) error {
	logger.Logger.Info("CreateCommonCategoryListService called", zap.Int("count", len(categories)))

	for i := range categories {
		categories[i].UserId = nil
		result := db.DB.Omit(clause.Associations).Create(&categories[i])
		if result.Error != nil {
			logger.Logger.Error("Database error in CreateCommonCategoryListService", zap.Error(result.Error))
			return result.Error
		}
	}
	logger.Logger.Info("Successfully created common category list")
	return nil
}

// DeleteCommonCategoryService 删除宠物行为的业务逻辑
func DeleteCommonCategoryService(id uint) error {
	logger.Logger.Info("DeleteCommonCategoryService called", zap.Uint("id", id))

	result := db.DB.Delete(&models.RecordCategory{}, "id = ? AND user_id IS NULL", id)
	if result.Error != nil {
		logger.Logger.Error("Database error in DeleteCommonCategoryService", zap.Error(result.Error))
		return result.Error
	}
	logger.Logger.Info("Successfully deleted common category", zap.Uint("id", id))
	return nil
}

// GetUserListService 获取用户列表的业务逻辑
func GetUserListService(page models.PageModel) ([]models.UserInfo, error) {
	logger.Logger.Info("GetUserListService called", zap.Int("pageNum", page.PageNum), zap.Int("pageSize", page.PageSize))

	var userList []models.UserInfo
	offer := (page.PageNum - 1) * page.PageSize

	result := db.DB.Model(models.UserInfo{}).
		Offset(offer).
		Limit(page.PageSize).
		Order("created_at DESC").
		Find(&userList)
	if result.Error != nil {
		logger.Logger.Error("Database error in GetUserListService", zap.Error(result.Error))
		return nil, result.Error
	}
	logger.Logger.Info("Successfully retrieved user list", zap.Int("count", len(userList)))
	return userList, nil
}

// GetLikeListService 点赞列表的业务逻辑
func GetLikeListService(page models.PageModel) ([]models.LikeMessageModel, error) {
	logger.Logger.Info("GetLikeListService called", zap.Int("pageNum", page.PageNum), zap.Int("pageSize", page.PageSize))

	var likeList []models.LikeMessageModel
	offset := (page.PageNum - 1) * page.PageSize
	result := db.DB.Model(models.LikeMessageModel{}).
		Offset(offset).
		Limit(page.PageSize).
		Order("created_at DESC").
		Find(&likeList)
	if result.Error != nil {
		logger.Logger.Error("Database error in GetLikeListService", zap.Error(result.Error))
		return nil, result.Error
	}
	logger.Logger.Info("Successfully retrieved like list", zap.Int("count", len(likeList)))
	return likeList, nil
}

// GetCollectionListService 收藏列表的业务逻辑
func GetCollectionListService(page models.PageModel) ([]models.CollectionMessageModel, error) {
	logger.Logger.Info("GetCollectionListService called", zap.Int("pageNum", page.PageNum), zap.Int("pageSize", page.PageSize))

	var collectionList []models.CollectionMessageModel
	offset := (page.PageNum - 1) * page.PageSize
	result := db.DB.Model(models.CollectionMessageModel{}).
		Offset(offset).Limit(page.PageSize).
		Order("created_at DESC").
		Find(&collectionList)

	if result.Error != nil {
		logger.Logger.Error("Database error in GetCollectionListService", zap.Error(result.Error))
		return nil, result.Error
	}
	logger.Logger.Info("Successfully retrieved collection list", zap.Int("count", len(collectionList)))
	return collectionList, nil
}