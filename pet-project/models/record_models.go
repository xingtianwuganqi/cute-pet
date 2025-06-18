package models

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt" `
}

// RecordList Type 是日常还花销 1：共同日常，2：自定义日常，3：共同花销，4：自定花销
// 是哪个宠物
type RecordList struct {
	BaseModel
	User            UserInfo          `json:"user"`
	UserId          uint              `json:"userId" form:"userId"`
	Type            uint              `json:"type" form:"type" gorm:"default:0"`
	PetInfo         PetInfo           `json:"petInfo"`
	PetInfoId       uint              `json:"petInfoId" form:"petInfoId"  binding:"required"`
	ActionType      *PetActionType    `json:"actionType" gorm:"-"`
	ActionTypeId    *uint             `json:"actionTypeId" form:"actionTypeId"`
	CustomAction    *PetCustomAction  `json:"customType" gorm:"-"`
	CustomActionId  *uint             `json:"customActionId" form:"customActionId"`
	ConsumeType     *PetConsumeType   `json:"consumeType" gorm:"-"`
	ConsumeTypeId   *uint             `json:"consumeTypeId" form:"consumeTypeId"`
	CustomConsume   *PetCustomConsume `json:"customConsume" gorm:"-"`
	CustomConsumeId *uint             `json:"customConsumeId" form:"customConsumeId"`
	Spend           float32           `json:"spend" form:"spend" gorm:"default:0"`
	Desc            string            `json:"desc" form:"desc" gorm:"default:''"`
}

// PetActionType 宠物日常/*
type PetActionType struct {
	BaseModel
	ActionName string `json:"actionName" form:"actionName" gorm:"size:64" binding:"required"`
	Icon       string `json:"icon" form:"icon" gorm:"size:64" binding:"required"`
}

// PetCustomAction  宠物日常自定义/*
type PetCustomAction struct {
	BaseModel
	User       UserInfo `json:"user" form:"user"`
	UserId     uint     `json:"userId" form:"userId"`
	CustomName string   `json:"customName" form:"customName" gorm:"size:32" binding:"required"`
	CustomIcon string   `json:"customIcon" form:"customIcon" gorm:"size:256" binding:"required"`
	Desc       string   `json:"desc" form:"desc" gorm:"size:256"`
}

// PetConsumeType 宠物消费
type PetConsumeType struct {
	BaseModel
	ConsumeName string  `json:"consumeName" form:"consumeName" gorm:"size:32" binding:"required"`
	ConsumeIcon float32 `json:"consumeIcon" form:"consumeIcon" gorm:"default:0" binding:"required"`
}

// PetCustomConsume 用户自定义消费
type PetCustomConsume struct {
	BaseModel
	User        UserInfo `json:"user" form:"user"`
	UserId      uint     `json:"userId" form:"userId"`
	ConsumeName string   `json:"consumeName" form:"consumeName" gorm:"size:32" binding:"required"`
	ConsumeIcon string   `json:"consumeIcon" form:"consumeIcon" gorm:"default:0" binding:"required"`
	Desc        string   `json:"desc" form:"desc" gorm:"size:256"`
}

// PetInfo
/*
pet_type : 0:默认值，1：猫咪，2：狗，3：其他
Unit: 1:kg 2:g 3:斤
Gender：1:公 2:母
*/
type PetInfo struct {
	BaseModel
	User     UserInfo `json:"user" form:"user"`
	UserId   uint     `json:"userId" form:"userId"`
	PetType  string   `json:"petType" form:"petType" gorm:"size:64"`
	Avatar   string   `json:"avatar" form:"avatar" gorm:"size:64" binding:"required"`
	Name     string   `json:"name" form:"name" gorm:"size:32" binding:"required"`
	Gender   uint     `json:"gender" form:"gender" gorm:"default:0"`
	BirthDay string   `json:"birthDay" form:"birthDay" gorm:"size:32"`
	HomeDay  string   `json:"homeDay" form:"homeDay" gorm:"size:32"`
	Desc     string   `json:"Desc" form:"Desc" gorm:"size:256"`
	Weight   float32  `json:"weight" form:"weight" gorm:"default:0"`
	Unit     uint     `json:"unit" form:"unit" gorm:"size:32"`
}

// AfterFind RecordList 查找其他字段
func (record *RecordList) AfterFind(tx *gorm.DB) (err error) {
	//if record.UserId != 0 {
	//	var user UserInfo
	//	result := tx.Model(&UserInfo{}).Where("user_id = ?", record.UserId).First(&user)
	//	if result.Error != nil {
	//		return result.Error
	//	} else {
	//		record.User = user
	//	}
	//}
	//if record.PetInfoId != 0 && record.UserId != 0 {
	//	var petInfo PetInfo
	//	result := tx.Model(&PetInfo{}).
	//		Where("id = ? AND user_id = ?", record.PetInfoId, record.UserId).
	//		First(&petInfo)
	//	if result.Error != nil {
	//		return result.Error
	//	} else {
	//		record.PetInfo = petInfo
	//	}
	//}
	if record.ActionTypeId != nil && record.UserId != 0 {
		var actionModel PetActionType
		result := tx.Model(&PetActionType{}).
			Where("action_type_id = ? AND user_id = ?", record.ActionTypeId, record.UserId).
			First(&actionModel)
		if result.Error != nil {
			return result.Error
		} else {
			record.ActionType = &actionModel
		}
	}

	if record.CustomActionId != nil && record.UserId != 0 {
		var customModel PetCustomAction
		result := tx.Model(&PetCustomAction{}).
			Where("custom_action_id = ? AND user_id = ?", record.CustomActionId, record.UserId).
			First(&customModel)
		if result.Error != nil {
			return result.Error
		} else {
			record.CustomAction = &customModel
		}
	}

	if record.ConsumeTypeId != nil && record.UserId != 0 {
		var consumeModel PetConsumeType
		result := tx.Model(&PetConsumeType{}).
			Where("consume_type_id = ? AND user_id = ?", record.ConsumeTypeId, record.UserId).
			First(&consumeModel)
		if result.Error != nil {
			return result.Error
		} else {
			record.ConsumeType = &consumeModel
		}
	}

	if record.CustomConsumeId != nil && record.UserId != 0 {
		var customConsumeModel PetCustomConsume
		result := tx.Model(&PetCustomConsume{}).
			Where("custom_consume_id = ? AND user_id = ?", record.CustomConsumeId, record.UserId).
			First(&customConsumeModel)
		if result.Error != nil {
			return result.Error
		} else {
			record.CustomConsume = &customConsumeModel
		}
	}

	return
}
