package handlers

import (
	"net/http"
	"log"

	"github.com/gorilla/websocket"
)

// TODO: Add ci & cd
// TODO: Add auth

const (
	readSize  = 1024
	writeSize = 1024

)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  readSize,
	WriteBufferSize: writeSize,
}

func HandleConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to websocket: %v", err)
		return
	}
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