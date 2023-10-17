package main

import "net/http"

func handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Client-side error: Something went wrong")
}
