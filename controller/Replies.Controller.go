package controller

import (
	"log"
	"net/http"
	"posthis/utils"
	"strconv"

	"github.com/gorilla/mux"
)

//Handlers
func GetReplies() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		replyModel := ReplyModel{Model: Model{Scheme: r.URL.Scheme, Host: r.URL.Host}}

		vars := mux.Vars(r)

		strId := vars["id"]

		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		models, err := replyModel.GetReplies(uint(id))
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: models, Message: "Replies retrieved sucessfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func CreateReply() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		replyModel := ReplyModel{}

		vars := mux.Vars(r)

		userId, err := strconv.ParseUint(vars["userId"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		postId, err := strconv.ParseUint(vars["postId"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.ParseMultipartForm(10 << 20)
		formdata := r.MultipartForm

		content := r.FormValue("content")
		files := formdata.File["files"]

		reply, err := replyModel.CreateReply(uint(userId), uint(postId), content, files)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: reply, Message: "Reply created successfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func UpdateReply() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		replyModel := ReplyModel{}

		vars := mux.Vars(r)
		strId := vars["id"]
		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//10mb total
		r.ParseMultipartForm(10 << 20)

		formdata := r.MultipartForm

		content := r.FormValue("content")
		files := formdata.File["files"]
		deleted := formdata.Value["deleted"]

		reply, err := replyModel.UpdateReply(uint(id), content, deleted, files)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: reply, Message: "Reply updated successfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func DeleteReply() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		replyModel := ReplyModel{}

		vars := mux.Vars(r)
		strId := vars["id"]

		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = replyModel.DeleteReply(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: nil, Message: "Reply deleted successfully"}

		utils.WriteJsonResponse(w, response)
	})
}
