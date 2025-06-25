package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"pet-project/models"
	"pet-project/settings"
)

var (
	ctx = context.Background()
	Rdb *redis.Client
	DB  *gorm.DB
	err error
)

func linkInit() {
	host := settings.Conf.Database.Host
	port := settings.Conf.Database.Port
	database := settings.Conf.Database.DataBase
	username := settings.Conf.Database.Username
	password := settings.Conf.Database.Password
	charset := settings.Conf.Database.Charset
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
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
	err := DB.AutoMigrate(
		&models.UserInfo{},
		&models.SuggestionModel{},

		&models.PetInfo{},
		&models.RecordCategory{},
		&models.CustomCategory{},
		&models.RecordList{},

		//&models.MessageModel{},
		//&models.LikeMessageModel{},
		//&models.CollectionMessageModel{},
		//&models.CommentModel{},
		//&models.ReplayModel{},
	)
	if err != nil {
		return
	}
}

func linkRedis() {
	addr := fmt.Sprintf("%s:%d", settings.Conf.Redis.Host, settings.Conf.Redis.Port)
	password := settings.Conf.Redis.Password
	redisDb := settings.Conf.Redis.DB
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       redisDb,
	})
	pong, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Redis connected to", pong)
}

func LinkDataBase() {
	linkInit()
	linkRedis()
}
