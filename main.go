package main

import (
	"log"
	"net/http"

	"github.com/ChipsAhoyEnjoyer/devWhisper/internal/server"
)

func main() {
	log.Println("devWhisper booting up...")

	config, err := server.NewServer()
	if err != nil {
		log.Fatalf("Failed to create a new server: %v", err)
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", config.AuthMiddleware(http.HandlerFunc(server.HandleHomepage)))
	mux.HandleFunc("/deleteUser", config.AuthMiddleware(http.HandlerFunc(config.HandleDeleteUser)))
	mux.HandleFunc("/connect", config.AuthMiddleware(http.HandlerFunc(server.HandleConnect)))
	mux.HandleFunc("/pingRedis", server.HandlePingRedis)
	mux.HandleFunc("/register", config.HandleRegister)
	mux.HandleFunc("/login", config.HandleLogin)
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}

	log.Printf("Serving devWhisper on port: %v\n", config.Port)
	log.Fatal(server.ListenAndServe())

}
