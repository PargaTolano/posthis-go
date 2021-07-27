package api

import (
	"encoding/json"
	"net/http"
	. "posthis/model"
	. "posthis/model/viewmodel"
	. "posthis/utils"

	"github.com/gorilla/context"
)

//Used Model: Post

func GetPosts(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Only GET method is acceptable for this request", http.StatusNotFound)
		return
	}

	db, err := ConnectToDb()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db.AutoMigrate(&Post{})

	var post Post

	db.Find(&post)

	marshal, err := json.Marshal(post)

	if err != nil {
		http.Error(w, "Only GET method is acceptable for this request", http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshal))
}

var CreatePost = http.HandlerFunc(createPost)

func createPost(w http.ResponseWriter, r *http.Request) {

	var (
		media    []*Media
		user     User
		post     Post
		response SuccesVM
	)

	//10mb total
	r.ParseMultipartForm(10 << 20)

	formdata := r.MultipartForm

	files := formdata.File["files"]

	err := UploadMultipleFiles(files, &media)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content := r.FormValue("content")

	db, err := ConnectToDb()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.AutoMigrate(&User{}, &Post{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	context.Get(r, "userId")

	db.First(&user, context.Get(r, "userId"))

	post = Post{Content: content, Media: media}

	db.Model(&user).Association("Posts").Append(&post)

	response = SuccesVM{Data: post, Message: "Post created successfully"}

	marshal, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(marshal)
}

var UpdatePost = http.HandlerFunc(updatePost)

func updatePost(w http.ResponseWriter, r *http.Request) {
}

var DeletePost = http.HandlerFunc(deletePost)

func deletePost(w http.ResponseWriter, r *http.Request) {
}

var GetFeed = http.HandlerFunc(getFeed)

func getFeed(w http.ResponseWriter, r *http.Request) {
}
