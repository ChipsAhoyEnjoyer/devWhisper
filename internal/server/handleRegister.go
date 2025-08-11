package server

import (
	"encoding/json"
	"github.com/ChipsAhoyEnjoyer/devWhisper/internal/database"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

func (cfg Config) HandleRegister(w http.ResponseWriter, r *http.Request) {
	//  TODO: Add input validation and password hashing
	log.Println("User is registering")
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	params := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Username:       params.Username,
		HashedPassword: params.Password,
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	})
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("User %s registered successfully with ID %s", user.Username, user.ID)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
