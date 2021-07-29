package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	. "posthis/db"
	. "posthis/model"
	. "posthis/model/viewmodel"
	. "posthis/utils"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

//Used Model: Post --main, User --owner-to-post, Media --belonging-to-post
//User ViewModels:

//Validation

//API handlers
func GetPosts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			posts    []Post
			response SuccesVM
		)

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.Find(&posts)

		response = SuccesVM{Data: posts, Message: "Retrieved Posts Sucessfully"}

		marshal, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(marshal))
	})
}

func GetPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			post         Post
			postDetailVm PostDetailVM
			response     SuccesVM
		)
		vars := mux.Vars(r)
		tid, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := uint(tid)

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.Preload("Media").First(&post, id)
		postDetailVm = PostDetailVM{ID: post.ID, Content: post.Content, Media: []string{}}

		for _, v := range post.Media {
			postDetailVm.Media = append(postDetailVm.Media, v.GetPath(r.URL.Scheme, r.Header.Get("Host")))
		}

		response = SuccesVM{Data: postDetailVm, Message: "Retrieved Posts Sucessfully"}

		marshal, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(marshal))
	})
}

//Takes its data from a multipart-form
//The key files holds the media to upload to the database and associate to this post
func CreatePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			media    []*Media
			user     User
			post     Post
			response SuccesVM
		)

		content := r.FormValue("content")

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.First(&user, context.Get(r, "userId"))

		//10mb total
		r.ParseMultipartForm(10 << 20)

		formdata := r.MultipartForm

		files := formdata.File["files"]

		err = UploadMultipleFiles(files, &media)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.CreateInBatches(&media, len(media))

		post = Post{Content: content, Media: media}

		db.Create(&post)
		db.Model(&user).Association("Posts").Append(&post)

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

//Takes its data from a multipart-form
//The key files holds the media to upload to the database and associate to this post
//The key json contains a json file that gives info on which media files are deleted
//and the new text content of the post
func UpdatePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			media    []*Media
			post     Post
			model    PostUpdateVM
			response SuccesVM
		)

		vars := mux.Vars(r)
		id := vars["id"]

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//10mb total
		r.ParseMultipartForm(10 << 20)

		formdata := r.MultipartForm

		files := formdata.File["files"]
		jsonData := formdata.File["json"]

		db.First(&post, id)

		if len(jsonData) == 1 {
			file, err := jsonData[0].Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			bytes, err := ioutil.ReadAll(file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = json.Unmarshal(bytes, &model)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(model.Deleted) > 0 {
				db.Delete(&Media{}, model.Deleted)
			}

			post.Content = model.Content

		} else if len(jsonData) > 1 {
			http.Error(w, "Can only receive one json data file", http.StatusBadRequest)
			return
		}

		if len(files) >= 1 {
			err = UploadMultipleFiles(files, &media)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			db.CreateInBatches(&media, len(media))
			db.Model(&post).Association("Media").Append(media)
		}

		db.Save(&post)

		response = SuccesVM{Data: post, Message: "Post updated successfully"}

		marshal, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(marshal)
	})
}

func DeletePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			post Post
		)

		vars := mux.Vars(r)
		id := vars["id"]

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tid, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		db.Delete(&post, uint(tid))

		marshal, err := json.Marshal(SuccesVM{Data: post, Message: "Post deleted successfully"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(marshal)
	})
}

func GetFeed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
