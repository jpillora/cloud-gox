package main

import (
	"log"
	"os"

	// "github.com/jpillora/cloud-gox/server"

	"./server/" //this warning is okay since we're 'go run'ning it
)

func main() {
	//run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	s := server.NewServer(port)
	log.Printf("listening on %s...", port)
	log.Fatal(s.Start())
}
