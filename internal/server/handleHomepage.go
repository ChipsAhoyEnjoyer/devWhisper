package server

import (
	"log"
	"net/http"
)

func HandleHomepage(w http.ResponseWriter, r *http.Request) {
	log.Println("User has landed in the homepage")
	_, err := w.Write([]byte("Welcome to devWhisper"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
