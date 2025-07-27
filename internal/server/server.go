package server

import (
	"os"
	"errors"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	maxConns = 100

	envFile = ".env"
)

var errEnvironmentVarNotSet = errors.New("error environment variable not set")
type activeConnections map[string]*websocket.Conn


type Config struct {
	Port string
	Users activeConnections

}


func NewServer() (*Config, error) {
	godotenv.Load(envFile)

	port := os.Getenv("PORT")
	if port == "" {
		return nil, errEnvironmentVarNotSet
	}

	conns := make(activeConnections, maxConns)

	c := Config{
		Port: port,
		Users: conns,
	}
	
	return &c, nil
}