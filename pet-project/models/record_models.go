package models

import (
	"gorm.io/gorm"
	"time"
)

// RecordList Type 是日常还花销
// 是哪个宠物
type RecordList struct {
	gorm.Model
	User              UserInfo             `json:"user" gorm:"foreignKey:UserId"`
	UserId            uint                 `json:"userId" form:"user_id"`
	Type              uint                 `json:"type" form:"type" gorm:"default:0"`
	PetInfo           PetInfo              `json:"petInfo"`
	PetInfoID         uint                 `json:"petInfoId" form:"petInfoId"`
	ActionType        PetActionType        `json:"action_type" gorm:"foreignKey:PetActionTypeId"`
	PetActionTypeId   uint                 `json:"petActionTypeId" form:"petActionTypeId"`
	CustomType        PetCustomType        `json:"custom_type" gorm:"foreignKey:PetCustomTypeId"`
	PetCustomTypeId   uint                 `json:"petCustomTypeId" form:"petCustomTypeId"`
	ConsumeType       PetConsumeType       `json:"consume_type" gorm:"foreignKey:ConsumeTypeId"`
	ConsumeTypeId     uint                 `json:"consume_type_id" form:"consumeTypeId"`
	CustomConsumeType PetCustomConsumeType `json:"custom_consume_type" gorm:"foreignKey:CustomConsumeId"`
	CustomConsumeId   uint                 `json:"custom_consume_id" form:"customConsumeId"`
	Spend             float32              `json:"spend" form:"spend" gorm:"default:0"`
	Desc              string               `json:"desc" form:"desc" gorm:"default:''"`
}

// PetActionType 宠物日常/*
type PetActionType struct {
	gorm.Model
	ActionName string `json:"action_name" form:"action_name"`
	Icon       string `json:"icon" form:"icon"`
}

// PetCustomType 宠物日常自定义/*
type PetCustomType struct {
	gorm.Model
	User       UserInfo `json:"user" form:"user" gorm:"foreignKey:UserId"`
	UserId     uint     `json:"userId" form:"userId"`
	CustomName string   `json:"customName" form:"customName" gorm:"size:32"`
	CustomIcon string   `json:"customIcon" form:"customIcon" gorm:"size:256"`
	Desc       string   `json:"desc" form:"desc" gorm:"size:256"`
}

type PetConsumeType struct {
	gorm.Model
	User        UserInfo `json:"user" form:"user" gorm:"foreignKey:UserId"`
	UserId      uint     `json:"userId" form:"userId"`
	ConsumeName string   `json:"consume_name" form:"consume_name" gorm:"size:32"`
	ConsumeIcon float32  `json:"consume_icon" form:"consume_icon" gorm:"default:0"`
	Desc        string   `json:"desc" form:"desc" gorm:"size:256"`
}

type PetCustomConsumeType struct {
	gorm.Model
	User        UserInfo `json:"user" form:"user" gorm:"foreignKey:UserId"`
	UserId      uint     `json:"userId" form:"userId"`
	ConsumeName string   `json:"consume_name" form:"consume_name" gorm:"size:32"`
	ConsumeIcon string   `json:"consume_icon" form:"consume_icon" gorm:"default:0"`
	Desc        string   `json:"desc" form:"desc" gorm:"size:256"`
}

// PetInfo
/*
pet_type : 0:默认值，1：猫咪，2：狗，3：其他
SteriliTime：绝育时间
*/
type PetInfo struct {
	gorm.Model
	User          UserInfo `json:"user" gorm:"foreignKey:UserId"`
	UserId        uint     `json:"userId" form:"userId"`
	PetType       uint     `json:"pet_type" form:"pet_type" gorm:"default:0"`
	Avatar        string   `json:"avatar" form:"avatar" gorm:"size:64"`
	Name          string   `json:"name" form:"name" gorm:"size:32"`
	Age           uint     `json:"age" form:"age" gorm:"default:0"`
	Gender        uint     `json:"gender" form:"gender" gorm:"default:0"`
	BirthDay      string   `json:"birthDay" form:"birthDay" gorm:"size:32"`
	HomeDay       time.Time
	Sterilization uint      `json:"sterilization" form:"sterilization" gorm:"default:0"`
	SteriliTime   time.Time `json:"steriliTime" form:"steriliTime" gorm:"default:0"`
	Desc          string    `json:"Desc" form:"Desc" gorm:"size:255"`
	Weight        float32   `json:"weight" form:"weight" gorm:"default:0"`
	Unit          string    `json:"unit" form:"unit" gorm:"size:32"`
}
