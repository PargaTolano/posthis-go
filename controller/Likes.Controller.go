package controller

import (
	"net/http"
	"posthis/utils"
	"strconv"

	"github.com/gorilla/mux"
)

//Handlers
func GetLikes() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		likeModel := LikeModel{}

		vars := mux.Vars(r)

		strId := vars["id"]

		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		likes, err := likeModel.GetLikes(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: likes, Message: "Replies retrieved sucessfully"}

		utils.WriteJsonResponse(w, response)
	})
}

//TODO prevent repetition
func CreateLike() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		likeModel := LikeModel{}

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

		like, err := likeModel.CreateLike(uint(userId), uint(postId))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: like, Message: "Like created successfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func DeleteLike() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		likeModel := LikeModel{}

		vars := mux.Vars(r)
		strId := vars["id"]

		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = likeModel.DeleteLike(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: nil, Message: "Like deleted successfully"}

		utils.WriteJsonResponse(w, response)
	})
}
