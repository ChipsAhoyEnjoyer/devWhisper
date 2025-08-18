package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, err error, response payload) {
	if code >= 500 {
		log.Println(response.Response)
		log.Println(err)
	}
	dat, err := json.Marshal(response)
	if err != nil {
		msg := "Failed to encode error response"
		log.Println(msg + ": " + err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func RespondWithJSON(w http.ResponseWriter, code int, response payload) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	dat, err := json.Marshal(response)
	if err != nil {
		msg := "Failed to encode response"
		log.Println(msg + ": " + err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Write(dat)
}
