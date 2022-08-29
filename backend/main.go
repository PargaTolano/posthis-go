package main

import (
	"log"
	"net/http"
	"os"

	"posthis/auth"
	"posthis/controller"
	"posthis/database"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"posthis/utils"
)

func initRoutes() http.Handler {

	r := mux.NewRouter()

	r.UseEncodedPath()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pages/index.html")
	})

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pages/404.html")
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", fileServer))

	r.Handle("/api/posts", auth.TokenAuthMiddleware(controller.GetPosts())).Methods("GET")
	r.Handle("/api/post/{id}", auth.TokenAuthMiddleware(controller.GetPost())).Methods("GET")
	r.Handle("/api/posts-create", auth.TokenAuthMiddleware(controller.CreatePost())).Methods("POST")
	r.Handle("/api/posts-update/{id}", auth.TokenAuthMiddleware(controller.UpdatePost())).Methods("PUT")
	r.Handle("/api/posts-delete/{id}", auth.TokenAuthMiddleware(controller.DeletePost())).Methods("DELETE")
	r.Handle("/api/posts-feed/{offset}/{limit}", auth.TokenAuthMiddleware(controller.GetFeed())).Methods("GET")
	r.Handle("/api/posts-feed/{id}/{offset}/{limit}", auth.TokenAuthMiddleware(controller.GetUserFeed())).Methods("GET")

	r.Handle("/api/users", auth.TokenAuthMiddleware(controller.GetUsers())).Methods("GET")
	r.Handle("/api/user/{id}", auth.TokenAuthMiddleware(controller.GetUser())).Methods("GET")
	r.Handle("/api/users-create", controller.CreateUser()).Methods("POST")
	r.Handle("/api/users-update/{id}", auth.TokenAuthMiddleware(controller.UpdateUser())).Methods("PUT")
	r.Handle("/api/users-delete/{id}", auth.TokenAuthMiddleware(controller.DeleteUser())).Methods("DELETE")
	r.Handle("/api/validate-password/{password}", auth.TokenAuthMiddleware(controller.ValidatePassword())).Methods("GET")
	r.Handle("/api/login", controller.Login()).Methods("POST")
	r.Handle("/api/logout", auth.TokenAuthMiddleware(controller.Logout())).Methods("POST")

	r.Handle("/api/search/{offset-post}/{limit-post}/{offset-user}/{limit-user}", auth.TokenAuthMiddleware(controller.GetSearch())).Methods("GET")

	r.Handle("/api/replies/{id}", auth.TokenAuthMiddleware(controller.GetReplies())).Methods("GET")
	r.Handle("/api/replies-create/{userId}/{postId}", auth.TokenAuthMiddleware(controller.CreateReply())).Methods("POST")
	r.Handle("/api/replies-update/{id}", auth.TokenAuthMiddleware(controller.UpdateReply())).Methods("PUT")
	r.Handle("/api/replies-delete/{id}", auth.TokenAuthMiddleware(controller.DeleteUser())).Methods("DELETE")

	r.Handle("/api/likes/{id}", auth.TokenAuthMiddleware(controller.GetLikes())).Methods("GET")
	r.Handle("/api/likes-create/{userId}/{postId}", auth.TokenAuthMiddleware(controller.CreateLike())).Methods("POST")
	r.Handle("/api/likes-delete/{id}", auth.TokenAuthMiddleware(controller.DeleteLike())).Methods("DELETE")

	r.Handle("/api/reposts/{id}", auth.TokenAuthMiddleware(controller.GetReposts())).Methods("GET")
	r.Handle("/api/reposts-create/{userId}/{postId}", auth.TokenAuthMiddleware(controller.CreateRepost())).Methods("POST")
	r.Handle("/api/reposts-delete/{id}", auth.TokenAuthMiddleware(controller.DeleteRepost())).Methods("DELETE")

	r.Handle("/api/follows/{id}", auth.TokenAuthMiddleware(controller.GetFollows())).Methods("GET")
	r.Handle("/api/follows-following/{id}", auth.TokenAuthMiddleware(controller.GetFollowing())).Methods("GET")
	r.Handle("/api/follows-create/{id}", auth.TokenAuthMiddleware(controller.CreateFollow())).Methods("POST")
	r.Handle("/api/follows-delete/{id}", auth.TokenAuthMiddleware(controller.DeleteFollow())).Methods("DELETE")

	http.Handle("/", r)

	return r
}

func main() {

	//Load enviroment variables
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}

	//Init auth redis Database
	auth.Init()

	//Init gorm
	database.InitDB()

	//Init firebase app
	utils.Init()

	mux := initRoutes()

	handler := cors.AllowAll().Handler(mux)

	port := os.Getenv("PORT")

	log.Println("Starting server on :" + port)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
