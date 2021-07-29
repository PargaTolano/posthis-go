package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"posthis/auth"
	. "posthis/db"
	. "posthis/model"
	. "posthis/model/viewmodel"
	. "posthis/utils"
	"strconv"

	"github.com/gorilla/mux"
)

//Used Model: User
//Used ViewModels: CreateUserVM,

//Validation
func validateCreateModel(model *UserCreateVM) error {
	return nil
}

func validateUpdateModel(model *UserCreateVM) error {
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
			model    UserCreateVM
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
			user     User
			model    UserUpdateVM
			response SuccesVM
		)

		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		db.First(&user, uint(id))
		if db.Error != nil {
			http.Error(w, db.Error.Error(), http.StatusInternalServerError)
		}

		r.ParseMultipartForm(10 << 20)

		formdata := r.MultipartForm

		profilePic := formdata.File["profilePic"]
		coverPic := formdata.File["coverPic"]
		jsonData := formdata.File["json"]

		if len(profilePic) == 1 {
			var media []*Media
			err = UploadMultipleFiles(profilePic, &media)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			db.Create(&media[0])
			user.ProfilePic = media[0]
		} else if len(profilePic) > 1 {
			http.Error(w, "Can only receive one profilePic data file", http.StatusBadRequest)
			return
		}

		if len(coverPic) == 1 {
			var media []*Media
			err = UploadMultipleFiles(coverPic, &media)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			db.Create(&media[0])
			user.CoverPic = media[0]
		} else if len(coverPic) > 1 {
			http.Error(w, "Can only receive one coverPic data file", http.StatusBadRequest)
			return
		}

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

			if model.Tag != "" {
				user.Tag = model.Tag
			}

			if model.Email != "" {
				user.Email = model.Email
			}

			if model.Username != "" {
				user.Username = model.Username
			}
		} else if len(coverPic) > 1 {
			http.Error(w, "Can only receive one json data file", http.StatusBadRequest)
			return
		}

		db.Save(&user)

		response = SuccesVM{Data: user, Message: "User updated successfully"}
		marshal, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(marshal)
	})
}

func DeleteUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		db.Delete(&User{}, uint(id))
		if db.Error != nil {
			http.Error(w, db.Error.Error(), http.StatusInternalServerError)
		}
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
