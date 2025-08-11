package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg Config) HandleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("User is attempting to log in")
	defer r.Body.Close()
	params := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	user, err := cfg.DB.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if user.HashedPassword != params.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Login successful",
	}
	w.Write([]byte(response.Message))
}
