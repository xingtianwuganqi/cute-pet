package service

import (
	"pet-project/db"
	"pet-project/models"

	"gorm.io/gorm/clause"
)

// GetCommonCategoriesService 获取宠物分类的业务逻辑
func GetCommonCategoriesService() ([]models.RecordCategory, error) {
	var petActionList []models.RecordCategory
	result := db.DB.Model(&models.RecordCategory{}).
		Where("user_id IS NULL").
		Find(&petActionList)
	if result.Error != nil {
		return nil, result.Error
	}
	return petActionList, nil
}

// CreateCommonCategoryService 创建宠物分类的业务逻辑
func CreateCommonCategoryService(recordCategory models.RecordCategory) error {
	recordCategory.UserId = nil
	result := db.DB.Create(&recordCategory)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateCommonCategoryListService 批量添加宠物行为的业务逻辑
func CreateCommonCategoryListService(categories []models.RecordCategory) error {
	for i := range categories {
		categories[i].UserId = nil
		result := db.DB.Omit(clause.Associations).Create(&categories[i])
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// DeleteCommonCategoryService 删除宠物行为的业务逻辑
func DeleteCommonCategoryService(id uint) error {
	result := db.DB.Delete(&models.RecordCategory{}, "id = ? AND user_id IS NULL", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetUserListService 获取用户列表的业务逻辑
func GetUserListService(page models.PageModel) ([]models.UserInfo, error) {
	var userList []models.UserInfo
	offer := (page.PageNum - 1) * page.PageSize

	result := db.DB.Model(models.UserInfo{}).
		Offset(offer).
		Limit(page.PageSize).
		Order("created_at DESC").
		Find(&userList)
	if result.Error != nil {
		return nil, result.Error
	}
	return userList, nil
}

// GetLikeListService 点赞列表的业务逻辑
func GetLikeListService(page models.PageModel) ([]models.LikeMessageModel, error) {
	var likeList []models.LikeMessageModel
	offset := (page.PageNum - 1) * page.PageSize
	result := db.DB.Model(models.LikeMessageModel{}).
		Offset(offset).
		Limit(page.PageSize).
		Order("created_at DESC").
		Find(&likeList)
	if result.Error != nil {
		return nil, result.Error
	}
	return likeList, nil
}

// GetCollectionListService 收藏列表的业务逻辑
func GetCollectionListService(page models.PageModel) ([]models.CollectionMessageModel, error) {
	var collectionList []models.CollectionMessageModel
	offset := (page.PageNum - 1) * page.PageSize
	result := db.DB.Model(models.CollectionMessageModel{}).
		Offset(offset).Limit(page.PageSize).
		Order("created_at DESC").
		Find(&collectionList)

	if result.Error != nil {
		return nil, result.Error
	}
	return collectionList, nil
}