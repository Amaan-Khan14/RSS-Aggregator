package main

import "net/http"

func hanlErr(w http.ResponseWriter, r *http.Request) {
	respondWithErr(w, 400, "Something went wrong")
}
