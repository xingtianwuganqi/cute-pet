package models

import (
	"gorm.io/gorm"
)

// RecordList Type 是日常还花销 1：共同日常，2：自定义日常，3：共同花销，4：自定花销
// 是哪个宠物
type RecordList struct {
	gorm.Model
	User               UserInfo             `json:"user" gorm:"-"`
	UserId             uint                 `json:"userId" form:"userId"`
	Type               uint                 `json:"type" form:"type" gorm:"default:0"`
	PetInfo            PetInfo              `json:"petInfo" gorm:"-"`
	PetId              uint                 `json:"petInfoId" form:"petInfoId"  binding:"required"`
	ActionModel        PetActionType        `json:"action_type" gorm:"-"`
	ActionId           uint                 `json:"petActionTypeId" form:"petActionTypeId" binding:"required"`
	CustomModel        PetCustomType        `json:"custom_type" gorm:"-"`
	CustomId           uint                 `json:"petCustomTypeId" form:"petCustomTypeId" binding:"required"`
	ConsumeModel       PetConsumeType       `json:"consume_type" gorm:"-"`
	ConsumeId          uint                 `json:"consume_type_id" form:"consumeTypeId" binding:"required"`
	CustomConsumeModel PetCustomConsumeType `json:"custom_consume_type" gorm:"-"`
	CustomConsumeId    uint                 `json:"custom_consume_id" form:"customConsumeId" binding:"required"`
	Spend              float32              `json:"spend" form:"spend" gorm:"default:0"`
	Desc               string               `json:"desc" form:"desc" gorm:"default:''"`
}

// PetActionType 宠物日常/*
type PetActionType struct {
	gorm.Model
	ActionName string `json:"action_name" form:"action_name" gorm:"size:64" binding:"required"`
	Icon       string `json:"icon" form:"icon" gorm:"size:64" binding:"required"`
}

// PetCustomType 宠物日常自定义/*
type PetCustomType struct {
	gorm.Model
	User       UserInfo `json:"user" form:"user" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	UserId     uint     `json:"userId" form:"userId"`
	CustomName string   `json:"customName" form:"customName" gorm:"size:32" binding:"required"`
	CustomIcon string   `json:"customIcon" form:"customIcon" gorm:"size:256" binding:"required"`
	Desc       string   `json:"desc" form:"desc" gorm:"size:256"`
}

// PetConsumeType 宠物消费
type PetConsumeType struct {
	gorm.Model
	ConsumeName string  `json:"consume_name" form:"consume_name" gorm:"size:32" binding:"required"`
	ConsumeIcon float32 `json:"consume_icon" form:"consume_icon" gorm:"default:0" binding:"required"`
}

// PetCustomConsumeType 用户自定义消费
type PetCustomConsumeType struct {
	gorm.Model
	User        UserInfo `json:"user" form:"user" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	UserId      uint     `json:"userId" form:"userId"`
	ConsumeName string   `json:"consume_name" form:"consume_name" gorm:"size:32" binding:"required"`
	ConsumeIcon string   `json:"consume_icon" form:"consume_icon" gorm:"default:0" binding:"required"`
	Desc        string   `json:"desc" form:"desc" gorm:"size:256"`
}

// PetInfo
/*
pet_type : 0:默认值，1：猫咪，2：狗，3：其他
Unit: 1:kg 2:g 3:斤
Gender：1:公 2:母
*/
type PetInfo struct {
	gorm.Model
	User     UserInfo `json:"user" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
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
	if record.UserId != 0 {
		var user UserInfo
		result := tx.Model(&UserInfo{}).Where("user_id = ?", record.UserId).First(&user)
		if result.Error != nil {
			return result.Error
		} else {
			record.User = user
		}
	}
	if record.PetId != 0 && record.UserId != 0 {
		var petInfo PetInfo
		result := tx.Model(&PetInfo{}).
			Where("id = ? AND user_id = ?", record.PetId, record.UserId).
			First(&petInfo)
		if result.Error != nil {
			return result.Error
		} else {
			record.PetInfo = petInfo
		}
	}
	if record.ActionId != 0 && record.UserId != 0 {
		var actionModel PetActionType
		result := tx.Model(&PetActionType{}).
			Where("action_id = ? AND user_id = ?", record.CustomId, record.UserId).
			First(&actionModel)
		if result.Error != nil {
			return result.Error
		} else {
			record.ActionModel = actionModel
		}
	}

	if record.CustomId != 0 && record.UserId != 0 {
		var customModel PetCustomType
		result := tx.Model(&PetCustomType{}).
			Where("custom_id = ? AND user_id = ?", record.CustomId, record.UserId).
			First(&customModel)
		if result.Error != nil {
			return result.Error
		} else {
			record.CustomModel = customModel
		}
	}

	if record.ConsumeId != 0 && record.UserId != 0 {
		var consumeModel PetConsumeType
		result := tx.Model(&PetConsumeType{}).
			Where("consume_id = ? AND user_id = ?", record.ConsumeId, record.UserId).
			First(&consumeModel)
		if result.Error != nil {
			return result.Error
		} else {
			record.ConsumeModel = consumeModel
		}
	}

	if record.CustomConsumeId != 0 && record.UserId != 0 {
		var customConsumeModel PetCustomConsumeType
		result := tx.Model(&PetCustomConsumeType{}).
			Where("custom_consume_id = ? AND user_id = ?", record.CustomConsumeId, record.UserId).
			First(&customConsumeModel)
		if result.Error != nil {
			return result.Error
		} else {
			record.CustomConsumeModel = customConsumeModel
		}
	}

	return
}
