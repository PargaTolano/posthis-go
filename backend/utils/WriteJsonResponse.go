package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter, data interface{}) {
	marshal, err := json.Marshal(data)

	if err != nil {
		log.Println("ERROR JSON ENCODING", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshal))
}
