package server

import "time"

const (
	refreshTokenDuration = time.Hour * 24 * 30
	jwtDuration          = time.Second * 15
)

type payload struct {
	Response string `json:"response"`
}

type token struct {
	RefreshToken string `json:"refresh_token"`
}
