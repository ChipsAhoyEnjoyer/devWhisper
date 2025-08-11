package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ChipsAhoyEnjoyer/devWhisper/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"time"
)

const (
	readSize  = 1024
	writeSize = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  readSize,
	WriteBufferSize: writeSize,
}

// TODO: Add authentication middleware & login handler

func (cfg Config) HandleRegister(w http.ResponseWriter, r *http.Request) {
	//  TODO: Add input validation and password hashing
	log.Println("User is registering")
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	params := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Username:       params.Username,
		HashedPassword: params.Password,
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	})
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("User %s registered successfully with ID %s", user.Username, user.ID)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func HandleHomepage(w http.ResponseWriter, r *http.Request) {
	log.Println("User has landed in the homepage")
	_, err := w.Write([]byte("Welcome to devWhisper"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (cfg Config) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	//  TODO: Validate incoming UUID
	//  TODO: Authenticate user before allowing deletion
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	params := struct {
		ID uuid.UUID `json:"id"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err = cfg.DB.DeleteUser(r.Context(), params.ID)
	if err != nil {
		http.Error(w, "Error deleting user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HandleConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to websocket: %v", err)
		return
	}
	go read(conn)
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
