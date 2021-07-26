package main

import (
	"log"
	"net/http"
	"posthis/api"
)

func routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/posts", api.GetPosts)

}

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
