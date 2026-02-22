package server

import (
	"context"
	"log"
	"net/http"
)

func (cfg Config) HandlePingRedis(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	resp, err := cfg.Rdb.Conn().Ping(ctx).Result()
	if err != nil {
		msg := payload{Response: "Failed to connect to redis"}
		RespondWithError(w, http.StatusInternalServerError, err, msg)
		return
	}

	msg := payload{Response: resp}
	RespondWithJSON(w, http.StatusOK, msg)
	log.Println("Connected to redis!")
}

func (cfg Config) HandleRedisSub(w http.ResponseWriter, r *http.Request) {
}
func (cfg Config) HandleRedisUnsub(w http.ResponseWriter, r *http.Request) {
}
func (cfg Config) HandleRedisPub(w http.ResponseWriter, r *http.Request) {
}
