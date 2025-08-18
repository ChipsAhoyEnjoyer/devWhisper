package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ChipsAhoyEnjoyer/devWhisper/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (cfg Config) HandleRegister(w http.ResponseWriter, r *http.Request) {
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

	// Username and password validation
	err = validateUsername(params.Username)
	if err != nil {
		msg := payload{Response: "invalid username\n" + err.Error()}
		RespondWithError(w, http.StatusBadRequest, err, msg)
		return
	}
	_, err = cfg.DB.GetUserByUsername(r.Context(), params.Username)
	if err == nil { // User exists
		msg := payload{Response: "username unavailable"}
		RespondWithError(w, http.StatusConflict, errors.New("username unavailable"), msg)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) { // Sql returns err if no row exists, so this checks for every other err
		msg := payload{Response: "internal error looking up username"}
		RespondWithError(w, http.StatusInternalServerError, err, msg)
		return

	}
	err = validatePassword(params.Password)
	if err != nil {
		msg := payload{Response: "invalid password\n" + err.Error()}
		RespondWithError(w, http.StatusBadRequest, err, msg)
		return
	}

	// Password hashing and storage
	hashPass, err := hashPassword(params.Password)
	if err != nil {
		msg := payload{Response: "server error storing password"}
		RespondWithError(w, http.StatusInternalServerError, err, msg)
		return
	}
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Username:       params.Username,
		HashedPassword: hashPass,
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	})
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	RespondWithJSON(
		w,
		http.StatusCreated,
		payload{
			Response: fmt.Sprintf("User %v created successfully at %v", user.Username, user.CreatedAt),
		},
	)
}
