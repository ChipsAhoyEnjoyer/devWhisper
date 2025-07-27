package main

import (
	"log"
	"net/http"

	"github.com/ChipsAhoyEnjoyer/devWhisper/internal/handlers"
	"github.com/ChipsAhoyEnjoyer/devWhisper/internal/server"
)

// const (
// 	successfulConnectionMsg = "New connection opened."
// )

func main() {
	log.Println("devWhisper booting up...")

	config, err := server.NewServer()
	if err != nil {
		log.Fatal()
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/connect", handlers.HandleConnect)

	server := http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}

	log.Printf("Serving devWhisper on port: %v\n", config.Port)
	log.Fatal(server.ListenAndServe())

}
