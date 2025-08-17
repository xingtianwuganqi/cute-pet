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

/*
	CategoryType 1.公共的，2，个人的
*/

type RecordCategory struct {
	BaseModel
	User         *UserInfo `json:"user" form:"user"`
	UserId       *uint     `json:"userId" form:"userId"`
	CategoryType uint      `json:"categoryType" form:"categoryType" gorm:"default:0"`
	Name         string    `json:"name" form:"name" binding:"required" gorm:"size:32" binding:"required"`
	Icon         string    `json:"icon" form:"icon" gorm:"size:64"`
	Color        string    `json:"color" form:"color" gorm:"size:32"`
	ImageUrl     string    `json:"imageUrl" form:"imageUrl" gorm:"size:64"`
	Desc         string    `json:"desc" form:"desc" gorm:"size:512"`
	Language     string    `json:"language" form:"language" gorm:"size:32"`
	Region       string    `json:"region" form:"region" gorm:"size:32"`
}

// PetInfo
/*
pet_type : 0:默认值，1：猫咪，2：狗，3：其他
Unit: 1:kg 2:g 3:斤
Gender：1:公 2:母
*/
type PetInfo struct {
	BaseModel
	User     *UserInfo `json:"user" form:"user"`
	UserId   uint      `json:"userId" form:"userId"`
	Avatar   string    `json:"avatar" form:"avatar" gorm:"size:64" binding:"required"`
	Name     string    `json:"name" form:"name" gorm:"size:32" binding:"required"`
	PetType  *string   `json:"petType" form:"petType" gorm:"size:32"`
	Gender   uint      `json:"gender" form:"gender" gorm:"default:0"`
	BirthDay string    `json:"birthDay" form:"birthDay" gorm:"size:32"`
	HomeDay  string    `json:"homeDay" form:"homeDay" gorm:"size:32"`
	Desc     string    `json:"desc" form:"desc" gorm:"size:256"`
	Weight   float32   `json:"weight" form:"weight" gorm:"default:0"`
	Unit     uint      `json:"unit" form:"unit" gorm:"size:32"`
	//Language string    `json:"language" form:"language" gorm:"size:32"`
	//Region   string    `json:"region" form:"region" gorm:"size:32"`
}

// RecordList Type 是日常还花销 1：共同日常，2：自定义日常，3：共同花销，4：自定花销
// 是哪个宠物
type RecordList struct {
	BaseModel
	User             *UserInfo       `json:"user"`
	UserId           uint            `json:"userId" form:"userId"`
	PetInfo          *PetInfo        `json:"petInfo" gorm:"-" binding:"-"`
	PetInfoId        uint            `json:"petInfoId" form:"petInfoId"`
	RecordCategory   *RecordCategory `json:"recordCategory" gorm:"-"`
	RecordCategoryId *uint           `json:"recordCategoryId" form:"recordCategoryId"`
	Spend            *float32        `json:"spend" form:"spend" gorm:"default:0"`
	Desc             string          `json:"desc" form:"desc" gorm:"size:512"  binding:"required"`
	Images           *StringArray    `json:"images" form:"images" gorm:"type:json"`
	RecordTime       time.Time       `json:"recordTime" form:"recordTime"`
	Language         string          `json:"language" form:"language" gorm:"size:32"`
	Region           string          `json:"region" form:"region" gorm:"size:32"`
}

// AfterFind RecordList 查找其他字段
func (record *RecordList) AfterFind(tx *gorm.DB) (err error) {
	if record.PetInfoId != 0 && record.UserId != 0 {
		var petInfo PetInfo
		result := tx.Model(&PetInfo{}).
			Omit("User").
			Where("id = ? AND user_id = ?", record.PetInfoId, record.UserId).
			First(&petInfo)
		if result.Error != nil {
			record.PetInfo = &PetInfo{}
		} else {
			record.PetInfo = &petInfo
		}
	}

	if record.RecordCategoryId != nil {
		var actionModel RecordCategory
		result := tx.Model(&RecordCategory{}).
			Where("id = ?", record.RecordCategoryId).
			First(&actionModel)
		if result.Error != nil {
			record.RecordCategory = nil
		} else {
			record.RecordCategory = &actionModel
		}
	}

	return
}
