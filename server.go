package main

import (
	"log"
	"net/http"
	"os"
)

func compile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/", compile)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
