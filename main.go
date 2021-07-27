package main

import (
	"log"
	"net/http"
	"posthis/api"
	"posthis/auth"

	"github.com/joho/godotenv"
)

func routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/posts", api.GetPosts)
	mux.Handle("/api/posts-create", auth.TokenAuthMiddleware(api.CreatePost))

	mux.HandleFunc("/api/users", api.GetUsers)
	mux.HandleFunc("/api/users-create", api.CreateUser)
	mux.Handle("/api/login", api.Login())
}

func main() {

	//Load enviroment variables
	godotenv.Load()

	auth.Init()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	routes(mux)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
