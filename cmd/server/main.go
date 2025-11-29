package main

import (
	"log"
	"os"

	"github.com/AmoghRao21/http-server-go/internal/server"
)

func main() {
	addr := ":8080"
	portEnv := os.Getenv("PORT")
	if portEnv != "" {
		addr = ":" + portEnv
	}

	srv := server.New(addr)
	err := srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
