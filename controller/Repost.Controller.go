package controller

import (
	"net/http"
	"posthis/utils"
	"strconv"

	"github.com/gorilla/mux"
)

//Handlers
func GetReposts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		repostModel := RepostModel{}

		vars := mux.Vars(r)

		strId := vars["id"]

		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reposts, err := repostModel.GetReposts(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: reposts, Message: "Reposts retrieved sucessfully"}

		utils.WriteJsonResponse(w, response)
	})
}

//TODO prevent repetition
func CreateRepost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		repostModel := RepostModel{}

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

		repost, err := repostModel.CreateRepost(uint(userId), uint(postId))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: repost, Message: "Repost created successfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func DeleteRepost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		repostModel := RepostModel{}

		vars := mux.Vars(r)
		strId := vars["id"]

		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = repostModel.DeleteRepost(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: nil, Message: "Repost deleted successfully"}

		utils.WriteJsonResponse(w, response)
	})
}
