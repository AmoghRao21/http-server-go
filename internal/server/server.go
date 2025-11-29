package server

import (
	"errors"
	"log"
	"net"
)

type Srv struct {
	Addr string
}

func New(addr string) *Srv {
	return &Srv{Addr: addr}
}

func (srv *Srv) Run() error {
	lstn, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	defer lstn.Close()

	log.Println("listening on", srv.Addr)

	for {
		conn, err := lstn.Accept()
		if err != nil {
			log.Println("accept err:", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		req, err := rdReq(conn)
		if err != nil {
			if errors.Is(err, errTooLarge) {
				writeResp(conn, 413, []byte("body too large"), "text/plain", false)
				log.Println("POST", "/data", 413, 0)
			}
			return
		}

		keep := wr(conn, req)
		if !keep {
			return
		}
	}
}
