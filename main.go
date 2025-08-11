package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/ChipsAhoyEnjoyer/devWhisper/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	envFile = ".env"

// successfulConnectionMsg = "New connection opened."
)

func main() {
	log.Println("devWhisper booting up...")

	godotenv.Load(envFile)
	db_url := os.Getenv("GOOSE_DBSTRING")
	if db_url == "" {
		log.Fatal("GOOSE_DBSTRING environment variable is not set")
		return
	}

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return
	}

	config, err := server.NewServer(db)
	if err != nil {
		log.Fatalf("Failed to create a new server: %v", err)
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.HandleHomepage)
	mux.HandleFunc("/register", config.HandleRegister)
	mux.HandleFunc("/deleteUser", config.HandleDeleteUser)
	mux.HandleFunc("/connect", server.HandleConnect)
	mux.HandleFunc("/login", config.HandleLogin)

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}

	log.Printf("Serving devWhisper on port: %v\n", config.Port)
	log.Fatal(server.ListenAndServe())

}
