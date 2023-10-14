package models

import (
	"gorm.io/gorm"
	"time"
)

type RecordList struct {
	gorm.Model
	ActionType      PetActionType `json:"action_type" gorm:"foreignKey:PetActionTypeId"`
	PetActionTypeId uint
	CustomType      PetCustomType `json:"custom_type" gorm:"foreignKey:PetCustomTypeId"`
	PetCustomTypeId uint
	Spend           float32
	Desc            string
}

type PetActionType struct {
	gorm.Model
	Type       uint `json:"type"`
	ActionName uint `json:"action_name"`
}

func (PetActionType) TableName() string {
	return "pet_action_type"
}

type PetCustomType struct {
	gorm.Model
	User   UserInfo `json:"user" gorm:"foreignKey:UserId"`
	UserId uint     `json:"userId"`
}

func (PetCustomType) TableName() string {
	return "pet_custom_info"
}

/*
pettype : 0:默认值，1：猫咪，2：狗，3：其他
*/

type PetInfo struct {
	gorm.Model
	User     UserInfo `json:"user" gorm:"foreignKey:UserId"`
	UserId   uint     `json:"userId"`
	PetType  uint     `json:"pet_type" gorm:"default:0"`
	Age      uint     `json:"age"`
	Name     string   `json:"name" gorm:"size:32"`
	BirthDay time.Time
}

func (PetInfo) TableName() string {
	return "pet_info"
}
