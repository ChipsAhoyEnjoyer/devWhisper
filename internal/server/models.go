package server

import (
	"time"

	"github.com/google/uuid"
)

const (
	refreshTokenDuration = time.Hour * 24 * 30
	jwtDuration          = time.Second * 15
)

type payload struct {
	Response string    `json:"response"`
	Id       uuid.UUID `json:"id"`
}
