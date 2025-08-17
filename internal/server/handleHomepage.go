package server

import (
	"log"
	"net/http"
)

func HandleHomepage(w http.ResponseWriter, r *http.Request) {
	log.Println("User has landed in the homepage")
	payload := struct {
		Message string `json:"response"`
	}{Message: "Welcome to devWhisper!"}
	RespondWithJSON(
		w,
		http.StatusOK,
		payload,
	)
}
