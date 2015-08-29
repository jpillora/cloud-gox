package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jpillora/cloud-gox/handler"
	"github.com/jpillora/opts"
)

var VERSION = "0.0.0-src"

func main() {
	//cli options
	var config = struct {
		Port int `help:"Listening port" env:"PORT"`
	}{
		Port: 3000,
	}
	//cli
	opts.New(&config).
		Name("cloud-gox").
		Version(VERSION).
		Repo("github.com/jpillora/cloud-gox").
		Parse()
	//run server
	h, err := handler.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listening on %d...", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), h))
}
