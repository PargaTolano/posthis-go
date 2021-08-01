package controller

import (
	"net/http"
	"strconv"

	"posthis/utils"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

//Handlers
func GetPosts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		postModel := PostModel{}

		posts, err := postModel.GetPosts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: posts, Message: "Retrieved Posts Sucessfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func GetPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		postModel := PostModel{Model: Model{Scheme: r.URL.Scheme, Host: r.URL.Host}}

		vars := mux.Vars(r)

		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		model, err := postModel.GetPost(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: model, Message: "Retrieved Post Sucessfully"}

		utils.WriteJsonResponse(w, response)
	})
}

//Takes its data from a multipart-form
//The key files holds the media to upload to the database and associate to this post
func CreatePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		postModel := PostModel{}

		r.ParseMultipartForm(10 << 20)

		formdata := r.MultipartForm

		ownerId := context.Get(r, "userId").(uint64)
		content := r.FormValue("content")
		files := formdata.File["files"]

		post, err := postModel.CreatePost(uint(ownerId), content, files)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: post, Message: "Post created successfully"}

		utils.WriteJsonResponse(w, response)
	})
}

//Takes its data from a multipart-form
//The key files holds the media to upload to the database and associate to this post
//The key json contains a json file that gives info on which media files are deleted
//and the new text content of the post
func UpdatePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		postModel := PostModel{}

		r.ParseMultipartForm(10 << 20)

		formdata := r.MultipartForm

		vars := mux.Vars(r)
		strId := vars["id"]

		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		post, err := postModel.UpdatePost(uint(id), r.FormValue("content"), formdata.Value["deleted"], formdata.File["files"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: post, Message: "Post updated successfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func DeletePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		postModel := PostModel{}

		vars := mux.Vars(r)
		strId := vars["id"]

		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = postModel.DeletePost(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: nil, Message: "Post deleted successfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func GetFeed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		postModel := PostModel{Model: Model{Scheme: r.URL.Scheme, Host: r.URL.Host}}

		vars := mux.Vars(r)
		offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id := context.Get(r, "userId").(uint64)

		models, err := postModel.GetFeed(uint(id), uint(offset), uint(limit))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: models, Message: "Retrieved Feed Sucessfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func GetUserFeed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		postModel := PostModel{Model: Model{Scheme: r.URL.Scheme, Host: r.URL.Host}}

		vars := mux.Vars(r)
		posterId, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		offset, err := strconv.ParseUint(vars["offset"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		limit, err := strconv.ParseUint(vars["limit"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id := context.Get(r, "userId").(uint64)

		models, err := postModel.GetUserFeed(uint(id), uint(posterId), uint(offset), uint(limit))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: models, Message: "Retrieved User Feed Sucessfully"}

		utils.WriteJsonResponse(w, response)
	})
}
