package server

import (
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	maxConns = 100

	envFile = ".env"

	defaultPort = "7777"
)

type activeConnections map[string]*websocket.Conn

type Config struct {
	Port string
	Users activeConnections

}


func NewServer() (*Config, error) {
	godotenv.Load(envFile)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	conns := make(activeConnections, maxConns)

	c := Config{
		Port: port,
		Users: conns,
	}
	
	return &c, nil
}