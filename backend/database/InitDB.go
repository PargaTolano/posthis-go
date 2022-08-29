package database

import (
	"log"
	"os"
	"posthis/entity"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User = entity.User
type Post = entity.Post
type Media = entity.Media
type Like = entity.Like
type Reply = entity.Reply
type Repost = entity.Repost
type Follow = entity.Follow

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NowFunc:                time.Now,
	})

	if err != nil {
		log.Println("Connection Failed to Open Exiting App...")
		os.Exit(-1)
	}

	log.Println("Connection Established here")

	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Could not retrieve sql connection...")
		os.Exit(-1)
	}

	sqlDB.SetMaxIdleConns(15)
	sqlDB.SetMaxOpenConns(151)

	DB = db

	//Add models to db
	DB.AutoMigrate(&User{}, &Post{}, &Media{}, &Like{}, &Reply{}, &Repost{}, &Follow{})
	DB.Raw("SET GLOBAL sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''));")
}
