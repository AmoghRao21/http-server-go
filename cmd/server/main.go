package main

import (
	"log"
	"os"

	"github.com/AmoghRao21/http-server-go/internal/server"
)

func main() {
	a := ":8080"
	if p := os.Getenv("Port"); p != "" {
		a = ":" + p
	}

	s := server.New(a)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
