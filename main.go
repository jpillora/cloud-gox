package main

import (
	"log"
	"os"

	"github.com/jpillora/cloud-gox/server"
)

func main() {
	//port
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	//run server
	s := server.NewServer(port)
	log.Printf("listening on %s...", port)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
