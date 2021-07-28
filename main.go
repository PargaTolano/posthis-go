package main

import (
	"log"
	"net/http"
	"posthis/api"
	"posthis/auth"
	"posthis/db"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func initRoutes() {

	r := mux.NewRouter()

	mux.CORSMethodMiddleware(r)

	r.UseEncodedPath()

	fileServer := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", fileServer))

	r.Handle("/api/posts", auth.TokenAuthMiddleware(api.GetPosts()))
	r.Handle("/api/post/{id}", api.GetPost())
	r.Handle("/api/posts-create", auth.TokenAuthMiddleware(api.CreatePost()))
	r.Handle("/api/posts-update/{id}", auth.TokenAuthMiddleware(api.UpdatePost()))
	r.Handle("/api/posts-delete/{id}", auth.TokenAuthMiddleware(api.DeletePost()))
	r.Handle("/api/posts-feed", auth.TokenAuthMiddleware(api.GetFeed()))

	r.Handle("/api/users", auth.TokenAuthMiddleware(api.GetUsers()))
	r.Handle("/api/users-create", api.CreateUser())
	r.Handle("/api/users-update/{id}", auth.TokenAuthMiddleware(api.UpdateUser()))
	r.Handle("/api/users-delete/{id}", auth.TokenAuthMiddleware(api.DeleteUser()))
	r.Handle("/api/login", api.Login())
	r.Handle("/api/logout", auth.TokenAuthMiddleware(api.Logout()))

	r.Handle("/api/search", auth.TokenAuthMiddleware(api.GetSearch()))

	r.Handle("/api/replies/{id}", auth.TokenAuthMiddleware(api.CreateReply()))
	r.Handle("/api/replies-create/{userId}/{postId}", auth.TokenAuthMiddleware(api.CreateReply()))
	r.Handle("/api/replies-update/{id}", auth.TokenAuthMiddleware(api.UpdateReply()))
	r.Handle("/api/replies-delete/{id}", auth.TokenAuthMiddleware(api.DeleteUser()))

	http.Handle("/", r)
}

func main() {

	//Load enviroment variables
	godotenv.Load()

	//Init auth redis Database
	auth.Init()

	//Init gorm
	db.InitDB()

	initRoutes()

	log.Println("Starting server on :4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
