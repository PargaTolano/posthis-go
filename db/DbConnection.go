package db

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToDb() (*gorm.DB, error) {
	dsn := "root:0000@tcp(127.0.0.1:3306)/posthis_local?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NowFunc:                time.Now,
	})
	return db, err
}
