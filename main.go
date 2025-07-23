package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	readSize  = 1024
	writeSize = 1024
	maxConns  = 100

	envFile                 = ".env"
	successfulConnectionMsg = "New connection opened."
)

type activeConnections map[string]*websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  readSize,
	WriteBufferSize: writeSize,
}

var errEnvironmentVarNotSet = errors.New("error environment variable not set")

func newConnectionMap() *activeConnections {
	new := make(activeConnections, maxConns)
	return &new
}

func (s *activeConnections) handleConnect(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to websocket: %v", err)
		return
	}
	log.Println(successfulConnectionMsg)

	read(conn)
}

func read(conn *websocket.Conn) {
	for {
		// read in a message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		log.Println(string(msg))
	}
}

func main() {
	log.Println("devWhisper booting up...")

	godotenv.Load(envFile)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("%v: %s", errEnvironmentVarNotSet, "PORT")
		return
	}

	conns := newConnectionMap()
	mux := http.NewServeMux()
	mux.HandleFunc("/connect", conns.handleConnect)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving devWhisper on port: %v\n", port)
	log.Fatal(server.ListenAndServe())

}
