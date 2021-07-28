package db

import . "posthis/model"

func InitDB() {
	db, err := ConnectToDb()
	if err != nil {
		panic("Error: " + err.Error())
	}

	//Add models to db
	db.AutoMigrate(&User{}, &Post{}, &Media{}, &Like{}, &Reply{}, &Repost{}, &Follow{})
}
