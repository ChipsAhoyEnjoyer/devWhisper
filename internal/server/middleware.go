package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type Middleware func(http.Handler) func(http.ResponseWriter, *http.Request)

func (cfg Config) AuthMiddleware(handler http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tk, err := GetBearerToken(r.Header)
		if err != nil {
			msg := payload{Response: "Unauthorized: not jwt found"}
			RespondWithError(w, http.StatusUnauthorized, err, msg)
			return
		}

		_, err = ValidateJWT(tk, cfg.TokenSecret)
		if err != nil {
			if !errors.Is(err, jwt.ErrTokenExpired) { // JWT wasnt valid
				msg := payload{Response: "error: unable to validate jwt; server error"}
				RespondWithError(w, http.StatusInternalServerError, err, msg)
				return
			} else { // Get new JWT if token is expired and Refresh token valid
				rTk := r.Header.Get("Refresh-Token")
				if rTk == "" {
					msg := payload{Response: "Unauthorized: Not logged in; No refresh token"}
					RespondWithError(w, http.StatusUnauthorized, errors.New("Unauthorized: Invalid refresh token"), msg)
					return
				}
				refreshTokenInfo, err := cfg.DB.GetRefreshTokenByToken(r.Context(), rTk)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						msg := payload{Response: "Unauthorized: Not logged in; No refresh token"}
						RespondWithError(w, http.StatusUnauthorized, errors.New("Unauthorized: Invalid refresh token"), msg)
						return
					}
					msg := payload{Response: "error confirming refresh token; server error"}
					RespondWithError(w, http.StatusInternalServerError, err, msg)
					return
				}
				jwtTk, err := MakeJWT(refreshTokenInfo.UserID, cfg.TokenSecret, jwtDuration)
				if err != nil {
					msg := payload{Response: "error creating jwt; server error"}
					RespondWithError(w, http.StatusInternalServerError, err, msg)
					return
				}
				w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtTk))
			}

		} // Either valid jwt and continue to endpoint, valid refresh token and get new jwt THEN go to new endpoint, or both invalid/missing and fail early
		handler.ServeHTTP(w, r)
	}
}
