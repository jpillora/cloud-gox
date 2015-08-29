package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jpillora/cloud-gox/handler"
)

func main() {
	//port
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	//run server
	h, err := handler.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listening on %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, h))
}
