package api

import (
	"encoding/json"
	"net/http"
	. "posthis/db"
	. "posthis/model"
	. "posthis/model/viewmodel"
	. "posthis/utils"
	"strconv"

	"github.com/gorilla/mux"
)

//Used Model: Post --main, User --owner-to-post, Media --belonging-to-post
//User ViewModels:

//Validation

//API handlers
func GetReplies() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			post     Post
			response SuccesVM
		)

		vars := mux.Vars(r)

		id := vars["id"]

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.AutoMigrate(&Post{}, &Reply{})

		db.First(&post, id)
		if db.Error != nil {
			http.Error(w, db.Error.Error(), http.StatusInternalServerError)
			return
		}

		response = SuccesVM{Data: post.Replies, Message: "Replies retrieved sucessfully"}

		marshal, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(marshal))
	})
}

func CreateReply() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			user     User
			post     Post
			reply    Reply
			media    []*Media
			response SuccesVM
		)

		vars := mux.Vars(r)

		tuid, err := strconv.ParseUint(vars["userId"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userId := uint(tuid)

		tpid, err := strconv.ParseUint(vars["postId"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		postId := uint(tpid)

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = db.AutoMigrate(&User{}, &Post{}, &Reply{}, &Media{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.First(&user, userId)
		if db.Error != nil {
			http.Error(w, db.Error.Error(), http.StatusInternalServerError)
			return
		}

		db.First(&post, postId)
		if db.Error != nil {
			http.Error(w, db.Error.Error(), http.StatusInternalServerError)
			return
		}

		//User and Post exist, so proceed to check media
		r.ParseMultipartForm(10 << 20)

		formdata := r.MultipartForm

		files := formdata.File["files"]

		err = UploadMultipleFiles(files, &media)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.CreateInBatches(&media, len(media))

		reply = Reply{Content: r.FormValue("content"), Media: media, UserID: userId, PostID: postId}

		//Once it works add everything to the database
		db.Create(&reply)
		db.Model(&user).Association("replies").Append(&reply)
		db.Model(&post).Association("replies").Append(&reply)

		response = SuccesVM{Data: post, Message: "Post created successfully"}

		marshal, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(marshal)
	})
}

func UpdateReply() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//vars := mux.Vars(r)
		//id := vars["id"]
	})
}

func DeleteReply() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//vars := mux.Vars(r)
		//id := vars["id"]
	})
}
