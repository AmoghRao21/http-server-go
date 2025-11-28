package server

import (
	"log"
	"net"
)

type S struct {
	A string
}

func New(a string) *S {
	return &S{A: a}
}

func (s *S) Run() error {
	ln, err := net.Listen("tcp", s.A)
	if err != nil {
		return err
	}

	defer ln.Close()
	log.Println("listening on", s.A)
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		c.Close()
	}
}
