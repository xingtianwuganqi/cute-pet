package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"pet-project/models"
)

var DB *gorm.DB

func LinkInit() {
	host := "localhost"
	port := "3306"
	database := "cutepet"
	username := "root"
	password := ""
	charset := "utf8"
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		username,
		password,
		host,
		port,
		database,
		charset)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error to DB connection ,err" + err.Error())
	}
	autoMigrateTable()
}

func autoMigrateTable() {
	DB.AutoMigrate(&models.UserInfo{}, &models.PetActionType{},
		&models.PetCustomType{}, &models.PetInfo{}, &models.RecordList{},
	)
}
