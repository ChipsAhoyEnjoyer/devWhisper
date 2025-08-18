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
	mux.HandleFunc("/", server.HandleHomepage)
	mux.HandleFunc("/register", config.HandleRegister)
	mux.HandleFunc("/deleteUser", config.HandleDeleteUser)
	mux.HandleFunc("/connect", server.HandleConnect)
	mux.HandleFunc("/login", config.HandleLogin)
	mux.HandleFunc("/ping", server.HandlePing)
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}

	log.Printf("Serving devWhisper on port: %v\n", config.Port)
	log.Fatal(server.ListenAndServe())

}
