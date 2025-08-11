package server

import (
	"log"
	"net/http"
)

func HandleHomepage(w http.ResponseWriter, r *http.Request) {
	log.Println("User has landed in the homepage")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to devWhisper"))
}
