package api

import "net/http"

func GetSearch(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Only GET method is acceptable for this request", http.StatusNotFound)
		return
	}

}
