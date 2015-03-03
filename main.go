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
	s, err := server.NewServer(port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listening on %s...", port)
	log.Fatal(s.Start())
}
