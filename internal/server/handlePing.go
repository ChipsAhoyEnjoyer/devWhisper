package server

import (
	"log"
	"net/http"

	"github.com/ChipsAhoyEnjoyer/devWhisper/internal/redis"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	addr := "localhost:6379"
	_, err := redis.NewRedisClient(addr)
	if err != nil {
		log.Printf("Failed to connect to redis: %v", err)
		return
	}
	log.Println("Connected to redis!")
}
