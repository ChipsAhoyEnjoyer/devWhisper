package server

import (
	"os"

	"github.com/ChipsAhoyEnjoyer/devWhisper/internal/database"
	"github.com/gorilla/websocket"
)

const (
	maxConns    = 100
	defaultPort = "7777"
)

type activeConnections map[string]*websocket.Conn

type Config struct {
	Port  string
	Users activeConnections
	DB    *database.Queries
}

func NewServer(db database.DBTX) (*Config, error) {

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	conns := make(activeConnections, maxConns)

	queries := database.New(db)
	c := Config{
		Port:  port,
		Users: conns,
		DB:    queries,
	}

	return &c, nil
}
