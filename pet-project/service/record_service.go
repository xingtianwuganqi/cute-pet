package service

import (
	"errors"
	"log"
	"pet-project/db"
	"pet-project/models"
	"pet-project/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PetInfoCreateService 提交宠物详情的业务逻辑
func PetInfoCreateService(userId uint, petInfo models.PetInfo) error {
	petInfo.UserId = userId
	result := db.DB.Omit(clause.Associations).Create(&petInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetPetListService 获取创建的宠物列表的业务逻辑
func GetPetListService(userId uint, pageModel models.PageModel) ([]models.PetInfo, error) {
	var petList []models.PetInfo
	offset := (pageModel.PageNum - 1) * pageModel.PageSize
	result := db.DB.
		Preload("User").
		Model(&models.PetInfo{}).
		Where("user_id = ?", userId).
		Offset(offset).
		Limit(pageModel.PageSize).
		Order("created_at DESC").
		Find(&petList)

	if result.Error != nil {
		return nil, result.Error
	}
	return petList, nil
}

// UpdatePetInfoService 更新宠物信息的业务逻辑
func UpdatePetInfoService(userId uint, petInfo models.PetInfo) error {
	if petInfo.ID == 0 {
		return errors.New("param error")
	}

	var oldPetInfo models.PetInfo
	oldResult := db.DB.Model(&models.PetInfo{}).Where("id = ?", petInfo.ID).First(&oldPetInfo)
	if errors.Is(oldResult.Error, gorm.ErrRecordNotFound) {
		return errors.New("data not exit")
	}
	if oldPetInfo.Avatar != petInfo.Avatar {
		// TODO: 删除旧头像
		// DeleteQiNiuFile(oldPetInfo.Avatar)
	}

	petInfo.UserId = userId
	result := db.DB.Model(&models.PetInfo{}).Where("id = ? AND user_id = ?", petInfo.ID, petInfo.UserId).
		Updates(models.PetInfo{
			UserId:   petInfo.UserId,
			PetType:  petInfo.PetType,
			Avatar:   petInfo.Avatar,
			Name:     petInfo.Name,
			Gender:   petInfo.Gender,
			BirthDay: petInfo.BirthDay,
			HomeDay:  petInfo.HomeDay,
			Weight:   petInfo.Weight,
			Unit:     petInfo.Unit,
			Desc:     petInfo.Desc,
		})
	if result.Error != nil {
		log.Println(result.Error)
		return errors.New("update error")
	}
	if result.RowsAffected == 0 {
		return errors.New("param error")
	}
	return nil
}

// DeletePetInfoService 删除创建的宠物详情的业务逻辑
func DeletePetInfoService(userId uint, petId uint) error {
	petInfo := models.PetInfo{}
	findResult := db.DB.Model(&models.PetInfo{}).
		Where("id = ? AND user_id = ?", petId, userId).
		First(&petInfo)
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		return errors.New("data not exit")
	}
	// TODO: 删除图片
	// DeleteQiNiuFile(petInfo.Avatar)
	result := db.DB.Delete(&petInfo, "id = ?", petId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetRecordCategoryListService 获取宠物行为列表的业务逻辑
func GetRecordCategoryListService(userId uint, pageModel models.CategoryTypeModel) ([]models.RecordCategory, error) {
	var petActionList []models.RecordCategory
	offset := (pageModel.PageNum - 1) * pageModel.PageSize
	if pageModel.CategoryType == 1 {
		result := db.DB.Model(&models.RecordCategory{}).
			Where("user_id = ?", userId).
			Offset(offset).
			Limit(pageModel.PageSize).
			Find(&petActionList)
		if result.Error != nil {
			return nil, result.Error
		}
		return petActionList, nil
	} else {
		result := db.DB.Model(&models.RecordCategory{}).
			Where("user_id IS NULL").Or("user_id = ?", userId).
			Offset(offset).
			Limit(pageModel.PageSize).
			Find(&petActionList)
		if result.Error != nil {
			return nil, result.Error
		}
		return petActionList, nil
	}
}

// CreateRecordCategoryService 添加宠物行为的业务逻辑
func CreateRecordCategoryService(userId uint, recordCategory models.RecordCategory) error {
	recordCategory.UserId = &userId
	result := db.DB.Create(&recordCategory)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateRecordCategoryService 更新宠物行为的业务逻辑
func UpdateRecordCategoryService(userId uint, recordCategory models.RecordCategory) error {
	var old models.RecordCategory
	if err := db.DB.Where("id = ? AND user_id = ?", recordCategory.ID, userId).First(&old).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("query error")
		}
		return err
	}

	if err := db.DB.Model(&old).Updates(recordCategory).Error; err != nil {
		return errors.New("update error")
	}

	return nil
}

// DeleteRecordCategoryService 删除宠物行为的业务逻辑
func DeleteRecordCategoryService(userId uint, id uint) error {
	recordCategory := models.RecordCategory{}
	findResult := db.DB.Model(&models.RecordCategory{}).Where("id = ? AND user_id = ?", id, userId).First(&recordCategory)
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		return errors.New("data not exit")
	}
	result := db.DB.Delete(&models.RecordCategory{}, "id = ? AND user_id = ?", id, userId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateRecordService 创建记录的业务逻辑
func CreateRecordService(userId uint, model models.RecordList) error {
	var petInfo models.PetInfo
	petResult := db.DB.Model(&models.PetInfo{}).Where("id = ?", model.PetInfoId).First(&petInfo)
	if errors.Is(petResult.Error, gorm.ErrRecordNotFound) {
		return errors.New("query error")
	}
	model.UserId = userId
	result := db.DB.Omit(clause.Associations).Create(&model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetRecordListService 查询记录列表的业务逻辑
func GetRecordListService(userId uint, pageModel models.RecordListModel) ([]models.RecordList, error) {
	var recordList []models.RecordList

	queryParam := models.RecordList{}
	if pageModel.CategoryId != nil {
		queryParam.RecordCategoryId = pageModel.CategoryId
	}
	if pageModel.PetInfoId != 0 {
		queryParam.PetInfoId = pageModel.PetInfoId
	}
	queryParam.UserId = userId

	offset := (pageModel.PageNum - 1) * pageModel.PageSize
	result := db.DB.Preload("User").
		Model(&models.RecordList{}).
		Where(&queryParam).
		Offset(offset).
		Limit(pageModel.PageSize).
		Order("record_time DESC").
		Find(&recordList)
	if result.Error != nil {
		return nil, result.Error
	}

	return recordList, nil
}

// DeleteRecordInfoService 删除记录的业务逻辑
func DeleteRecordInfoService(userId uint, id uint) error {

	record := models.RecordList{}
	findResult := db.DB.Model(&models.RecordList{}).
		Where("id = ? AND user_id = ?", id, userId).
		First(&record)

	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		return errors.New("data not exit")
	}

	// TODO: 删除图片
	if record.Images != nil {
		for _, image := range *record.Images {
			util.DeleteQiNiuFile(image)
		}
	}

	result := db.DB.Delete(&record)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetPetCostListService 查询宠物花费列表的业务逻辑
func GetPetCostListService(userId uint, pageModel models.RecordListModel) ([]models.RecordList, error) {
	var recordList []models.RecordList

	offset := (pageModel.PageNum - 1) * pageModel.PageSize

	tx := db.DB.
		Model(&models.RecordList{}).
		Where("user_id = ? AND spend > ?", userId, 0)

	if pageModel.PetInfoId != 0 {
		tx = tx.Where("pet_info_id = ?", pageModel.PetInfoId)
	}

	if pageModel.CategoryId != nil {
		tx = tx.Where("category_id = ?", *pageModel.CategoryId)
	}

	if pageModel.StartTime != nil {
		tx = tx.Where("record_time >= ?", *pageModel.StartTime)
	}

	if pageModel.EndTime != nil {
		tx = tx.Where("record_time <= ?", *pageModel.EndTime)
	}

	result := tx.
		Order("record_time DESC").
		Offset(offset).
		Limit(pageModel.PageSize).
		Find(&recordList)

	if result.Error != nil {
		return nil, result.Error
	}

	return recordList, nil
}
