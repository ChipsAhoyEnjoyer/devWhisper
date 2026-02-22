package server

import (
	"time"

	"github.com/ChipsAhoyEnjoyer/devWhisper/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const (
	refreshTokenDuration = time.Hour * 24 * 30
	jwtDuration          = time.Second * 15
)

type payload struct {
	Response string    `json:"response"`
	Id       uuid.UUID `json:"id"`
}

type redisChannel string

type activeConnections map[redisChannel][]*websocket.Conn

type Config struct {
	TokenSecret string
	Port        string
	Users       activeConnections
	Rdb         *redis.Client
	DB          *database.Queries
}
