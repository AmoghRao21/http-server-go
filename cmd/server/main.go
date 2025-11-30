package main

import (
	"flag"
	"log"
	"os"

	"github.com/AmoghRao21/http-server-go/internal/server"
)

func main() {
	loadtest := flag.Bool("load", false, "disable logging for load tests")
	flag.Parse()
	server.LoadTest = *loadtest

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
