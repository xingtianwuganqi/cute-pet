package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"pet-project/logger"
	"pet-project/models"
	"pet-project/settings"
	"go.uber.org/zap"
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
		logger.Logger.Error("Database connection failed", zap.Error(err))
		panic("Error to DB connection ,err" + err.Error())
	}
	logger.Logger.Info("Database connected successfully", zap.String("database", database), zap.String("host", host))
	autoMigrateTable()
}

func autoMigrateTable() {
	logger.Logger.Info("Starting auto migration of tables")
	err := DB.AutoMigrate(
		&models.UserInfo{},
		&models.SuggestionModel{},

		&models.PetInfo{},
		&models.RecordCategory{},
		&models.RecordList{},

		&models.TopicModel{},
		&models.PostModel{},
		&models.LikeMessageModel{},
		&models.CollectionMessageModel{},
		&models.MessageModel{},
		&models.CommentModel{},
		&models.ReplyModel{},
	)
	if err != nil {
		logger.Logger.Error("Auto migration failed", zap.Error(err))
		return
	}
	logger.Logger.Info("Auto migration completed successfully")
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
		logger.Logger.Error("Could not connect to Redis", zap.Error(err))
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	logger.Logger.Info("Redis connected successfully", zap.String("addr", addr), zap.String("result", pong))
}

func LinkDataBase() {
	logger.Logger.Info("Starting database connection process")
	linkInit()
	linkRedis()
	logger.Logger.Info("Database connection process completed")
}