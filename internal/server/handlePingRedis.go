package server

import (
	"log"
	"net/http"

	"github.com/ChipsAhoyEnjoyer/devWhisper/internal/redis"
)

func HandlePingRedis(w http.ResponseWriter, r *http.Request) {
	addr := "localhost:6379"
	_, err := redis.NewRedisClient(addr)
	if err != nil {
		msg := payload{Response: "Failed to connect to redis"}
		RespondWithError(w, http.StatusInternalServerError, err, msg)
		return
	}
	msg := payload{Response: "OK"}
	RespondWithJSON(w, http.StatusOK, msg)
	log.Println("Connected to redis!")
}
