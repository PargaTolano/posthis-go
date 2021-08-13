package controller

import (
	"net/http"

	"posthis/utils"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

//Handlers
func GetFollows() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		followModel := FollowModel{}

		vars := mux.Vars(r)

		strId := vars["id"]

		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		users, err := followModel.GetFollows(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: users, Message: "Followers retrieved sucessfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func GetFollowing() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		followModel := FollowModel{}

		vars := mux.Vars(r)

		strId := vars["id"]
		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		users, err := followModel.GetFollowing(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: users, Message: "Following retrieved sucessfully"}

		utils.WriteJsonResponse(w, response)
	})
}

//TODO make sure Follow doesn't exist already
func CreateFollow() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		followModel := FollowModel{Model: Model{Scheme: r.URL.Scheme, Host: r.URL.Host}}

		vars := mux.Vars(r)

		strFollowedId := vars["id"]

		followedId, err := strconv.ParseUint(strFollowedId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		followerId := context.Get(r, "userId").(uint64)

		user, err := followModel.CreateFollow(uint(followedId), uint(followerId))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: user, Message: "Follow created successfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func DeleteFollow() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		followModel := FollowModel{Model: Model{Scheme: r.URL.Scheme, Host: r.URL.Host}}

		vars := mux.Vars(r)
		strId := vars["id"]

		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		followerId := context.Get(r, "userId").(uint64)

		model, err := followModel.DeleteFollow(uint(followerId), uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: model, Message: "Follow deleted successfully"}

		utils.WriteJsonResponse(w, response)
	})
}
