package server

import (
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
		go func(conn net.Conn) {
			defer conn.Close()
			req, perr := rdReq(conn)
			if perr != nil {
				return
			}
			wr(conn, req)
		}(conn)
	}
}
