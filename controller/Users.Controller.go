package controller

import (
	"encoding/json"
	"net/http"
	"posthis/utils"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

//Handlers
func GetUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userModel := UserModel{}

		users, err := userModel.GetUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: users, Message: "Users retrieved Succesfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func GetUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userModel := UserModel{Model: Model{Scheme: r.URL.Scheme, Host: r.URL.Host}}

		vars := mux.Vars(r)

		strId := vars["id"]
		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		viewerId := context.Get(r, "userId").(uint64)

		user, err := userModel.GetUser(uint(id), uint(viewerId))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: user, Message: "User retrieved Succesfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func CreateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userModel := UserModel{}

		model := UserCreateVM{}

		err := json.NewDecoder(r.Body).Decode(&model)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := userModel.CreateUser(model)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := SuccesVM{Data: user, Message: "User created successfully"}

		utils.WriteJsonResponse(w, response)
	})
}

func UpdateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userModel := UserModel{}

		vars := mux.Vars(r)

		strId := vars["id"]
		id, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.ParseMultipartForm(10 << 20)

		formdata := r.MultipartForm

		pfpFiles := formdata.File["profilePic"]
		coverFiles := formdata.File["coverPic"]

		if len(pfpFiles) > 1 {
			http.Error(w, "Can only receive one file for profile picture", http.StatusBadRequest)
			return
		}

		if len(coverFiles) > 1 {
			http.Error(w, "Can only receive one file for cover picture", http.StatusBadRequest)
			return
		}

		model := UserUpdateVM{Tag: r.FormValue("tag"), Email: r.FormValue("email"), Username: r.FormValue("username"), Password: r.FormValue("password")}

		user, err := userModel.UpdateUser(uint(id), model, pfpFiles, coverFiles)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: user, Message: "User updated successfully"}
		utils.WriteJsonResponse(w, response)
	})
}

func DeleteUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userModel := UserModel{}

		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = userModel.DeleteUser(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: nil, Message: "User Deleted Succesfully!"}

		utils.WriteJsonResponse(w, response)
	})
}

func Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userModel := UserModel{}
		model := UserLoginVM{}

		err := json.NewDecoder(r.Body).Decode(&model)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tokens, err := userModel.Login(model)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: tokens, Message: "Login Succesful"}
		utils.WriteJsonResponse(w, response)
	})
}

func Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userModel := UserModel{}

		err := userModel.Logout(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: nil, Message: "Logged out succesfully"}

		utils.WriteJsonResponse(w, response)
	})
}
