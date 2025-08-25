package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ChipsAhoyEnjoyer/devWhisper/internal/database"
)

func (cfg Config) HandleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Authenticate user
	params := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		msg := payload{Response: "error: unable to read user request; server error"}
		RespondWithError(w, http.StatusInternalServerError, err, msg)
		return
	}
	user, err := cfg.DB.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msg := payload{Response: "error: invalid credentials"}
			RespondWithError(w, http.StatusUnauthorized, errors.New("Invalid credentials"), msg)
			return
		}
		msg := payload{Response: "error: unable to find user at this moment; server error"}
		RespondWithError(w, http.StatusInternalServerError, err, msg)
		return
	}
	err = CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		msg := payload{Response: "error: invalid credentials"}
		RespondWithError(w, http.StatusUnauthorized, errors.New("Invalid credentials"), msg)
		return
	}
	// Access tokens
	jwtTk, err := MakeJWT(user.ID, cfg.TokenSecret, jwtDuration)
	if err != nil {
		msg := payload{Response: "error: creating jwt; server error"}
		RespondWithError(w, http.StatusInternalServerError, err, msg)
		return
	}
	refreshToken, err := MakeRefreshToken()
	if err != nil {
		msg := payload{Response: "error: creating refresh token; server error"}
		RespondWithError(w, http.StatusInternalServerError, err, msg)
		return
	}
	_, err = cfg.DB.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:     refreshToken,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(refreshTokenDuration),
			RevokedAt: sql.NullTime{Valid: false},
			UserID:    user.ID,
		},
	)
	if err != nil {
		msg := payload{Response: "error: server could not save your refresh token; server error"}
		RespondWithError(w, http.StatusInternalServerError, err, msg)
		return
	}
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtTk))
	w.Header().Add("Refresh-Token", refreshToken)
	msg := payload{
		Response: fmt.Sprintf("Logging in as %s", user.Username),
		Id:       user.ID,
	}
	RespondWithJSON(w, http.StatusOK, msg)
}
