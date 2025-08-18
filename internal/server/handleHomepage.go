package server

import (
	"log"
	"net/http"
)

func HandleHomepage(w http.ResponseWriter, r *http.Request) {
	log.Println("User has landed in the homepage")
	msg := payload{Response: "Welcome to devWhisper!"}
	RespondWithJSON(
		w,
		http.StatusOK,
		msg,
	)
}
