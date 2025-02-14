package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dpatrie/sparkgrid/services"
)

func main() {
	s1, err := services.NewS1()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5555"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), s1))
}
