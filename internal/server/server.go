package server

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"net"
	"syscall"
	"time"
)

type Srv struct {
	Addr string
}

func New(addr string) *Srv {
	return &Srv{Addr: addr}
}

func (srv *Srv) Run() error {
	cfg := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var err error
			c.Control(func(fd uintptr) {
				applySocketOptions(fd)
			})
			return err
		},
	}

	lstn, err := cfg.Listen(context.Background(), "tcp", srv.Addr)
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
		if tc, ok := conn.(*net.TCPConn); ok {
			tc.SetNoDelay(true)
			tc.SetKeepAlive(true)
			tc.SetKeepAlivePeriod(30 * time.Second)
			tc.SetReadBuffer(1 << 20)
			tc.SetWriteBuffer(1 << 20)
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	br := bufio.NewReaderSize(conn, 64*1024)

	for {
		conn.SetReadDeadline(time.Now().Add(15 * time.Second))

		req, err := rdReq(br)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				return
			}
			if errors.Is(err, errTooLarge) {
				conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
				writeResp(conn, 413, []byte("body too large"), "text/plain", false)
				return
			}
			return
		}

		keep := shouldKeepAlive(req)

		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		keep = wr(conn, req) && keep

		if !keep {
			return
		}
	}
}
