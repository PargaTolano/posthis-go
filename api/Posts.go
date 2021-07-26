package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	. "posthis/model"
	. "posthis/utils"
)

//Used Model: Post

func GetPosts(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		fmt.Fprint(w, "400 Bad Request")
	}

	db := ConnectToDb()

	db.AutoMigrate(&Post{})

	var post Post

	db.Find(&post)

	w.Header().Add("Content-Type", "application/json")

	marshal, err := json.Marshal(post)

	if err != nil {
		panic("There was an error obtaining the posts : " + err.Error())
	}

	w.Write([]byte(marshal))
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
}
