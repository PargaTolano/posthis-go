package api

import (
	"encoding/json"
	"net/http"
	"posthis/auth"
	. "posthis/db"
	. "posthis/model"
	. "posthis/model/viewmodel"
	. "posthis/utils"

	"github.com/gorilla/mux"
)

//Used Model: User
//Used ViewModels: CreateUserVM,

//Validation
func validateCreateModel(model *CreateUserVM) error {
	return nil
}

func validateUpdateModel(model *CreateUserVM) error {
	return nil
}

//API handlers
func GetUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			users    []User
			response SuccesVM
		)

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.AutoMigrate(&User{})

		db.Find(&users)

		response = SuccesVM{Data: users, Message: "Users retrieved Succesfully"}

		marshal, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(marshal))
	})
}

func CreateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			model    CreateUserVM
			user     User
			response SuccesVM
		)

		err := json.NewDecoder(r.Body).Decode(&model)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = validateCreateModel(&model)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hash, err := HashPassword(model.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.AutoMigrate(&User{})

		user = User{Tag: model.Tag, Email: model.Email, Username: model.Username, PasswordHash: hash}
		db.Create(&user)

		response = SuccesVM{Data: user, Message: "User created successfully"}
		marshal, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(marshal)
	})
}

func UpdateUser() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			user User
		)

		vars := mux.Vars(r)

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		db.AutoMigrate(&User{})
		db.First(&user, vars["id"])
		if db.Error != nil {
			http.Error(w, db.Error.Error(), http.StatusInternalServerError)
		}

	})
}

func DeleteUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			model    LoginUserVM
			user     User
			response SuccesVM
		)

		err := json.NewDecoder(r.Body).Decode(&model)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := db.Where("username = ?", model.Username).First(&user).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !CheckPasswordHash(model.Password, user.PasswordHash) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ts, err := auth.CreateToken(uint64(user.ID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		saveErr := auth.CreateAuth(uint64(user.ID), ts)
		if saveErr != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}

		response = SuccesVM{Data: tokens, Message: "Login Succesful"}
		marshal, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(marshal)
	})
}

func Logout() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		au, err := auth.ExtractTokenMetadata(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		deleted, delErr := auth.DeleteAuth(au.AccessUuid)
		if delErr != nil || deleted == 0 {
			http.Error(w, delErr.Error(), http.StatusUnauthorized)
			return
		}

		marshal, err := json.Marshal(SuccesVM{Data: nil, Message: "Logget out succesfully"})

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(marshal))
	})
}
