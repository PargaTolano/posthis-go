package api

import (
	"encoding/json"
	"net/http"
	"posthis/auth"
	. "posthis/model"
	. "posthis/model/viewmodel"
	. "posthis/utils"
)

//Used Model: User
//Used ViewModels: CreateUserVM,

//Validation
func validateModel(model *CreateUserVM) error {
	return nil
	//return errors.New("Validation didnt work")
}

//API handlers
func GetUsers(w http.ResponseWriter, r *http.Request) {

	var (
		users    []User
		response SuccesVM
	)

	if r.Method != "GET" {
		http.Error(w, "Only GET method is acceptable for this request", http.StatusNotFound)
		return
	}

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
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

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

	err = validateModel(&model)

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
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
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
