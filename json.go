package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithErr(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responidng with 5XX error:", msg)
	}

	type errResponse struct {
		Error string `json:"erorr"` //{ "error": "message" }
	}

	respondWithJSON(w, code, errResponse{Error: msg})

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Set the content type to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("FAiled to Marshal JSON respond %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
