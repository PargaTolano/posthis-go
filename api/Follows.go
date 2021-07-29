package api

import (
	"encoding/json"
	"net/http"
	. "posthis/db"
	. "posthis/model"
	. "posthis/model/viewmodel"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

//Used Model: Follow --main, User --follower --followed
//User ViewModels:

//Validation

//API handlers
func GetFollows() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			users    []User
			response SuccesVM
		)

		vars := mux.Vars(r)

		id := vars["id"]

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.Joins("JOIN follows ON follows.follower_id = users.id AND follows.followed_id = ?", id).Find(&users)
		if db.Error != nil {
			http.Error(w, db.Error.Error(), http.StatusInternalServerError)
			return
		}

		response = SuccesVM{Data: users, Message: "Followers retrieved sucessfully"}

		marshal, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(marshal))
	})
}

func GetFollowing() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			users    []User
			response SuccesVM
		)

		vars := mux.Vars(r)

		id := vars["id"]

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.Joins("JOIN follows ON follows.followed_id = users.id AND follows.follower_id = ?", id).Find(&users)
		if db.Error != nil {
			http.Error(w, db.Error.Error(), http.StatusInternalServerError)
			return
		}

		response = SuccesVM{Data: users, Message: "Following retrieved sucessfully"}

		marshal, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(marshal))
	})
}

//TODO make sure Follow doesn't exist already
func CreateFollow() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			follow   Follow
			user     User
			follower User
			response SuccesVM
		)

		vars := mux.Vars(r)

		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		followerId := context.Get(r, "userId")

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.First(&user, id)
		db.First(&follower, followerId)

		if db.Error != nil {
			http.Error(w, db.Error.Error(), http.StatusInternalServerError)
			return
		}

		follow = Follow{FollowerID: follower.ID, FollowedID: user.ID}
		db.Create(&follow)
		db.Model(&user).Association("followers").Append(&follow)
		db.Model(&follower).Association("followings").Append(&follow)

		response = SuccesVM{Data: user, Message: "Follow created successfully"}

		marshal, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(marshal)
	})
}

func DeleteFollow() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			user     User
			follower User
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

		followerId := context.Get(r, "userId")

		db.First(&user, tid)
		db.First(&follower, followerId)

		db.Model(&user).Association("Followings").Delete(&follower)
		db.Model(&follower).Association("Followers").Delete(&user)

		marshal, err := json.Marshal(SuccesVM{Data: user, Message: "Post deleted successfully"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(marshal)
	})
}
